package stream

import (
   "41.neocities.org/dash"
   "41.neocities.org/sofia/file"
   "41.neocities.org/sofia/pssh"
   "bufio"
   "bytes"
   "encoding/base64"
   "errors"
   "fmt"
   "io"
   "log"
   "net/http"
   "net/url"
   "os"
   "slices"
   "strings"
   "time"
)

var ThreadCount = 1

func Mpd(name string, resp *http.Response) error {
   data, err := marshal(resp)
   if err != nil {
      return err
   }
   log.Println("WriteFile", name)
   err = os.WriteFile(name, data, os.ModePerm)
   if err != nil {
      return err
   }
   resp, err = unmarshal(data)
   if err != nil {
      return err
   }
   defer resp.Body.Close()
   data, err = io.ReadAll(resp.Body)
   if err != nil {
      return err
   }
   var mpd1 dash.Mpd
   err = mpd1.Unmarshal(data)
   if err != nil {
      return err
   }
   represents := slices.SortedFunc(mpd1.Representation(),
      func(a, b dash.Representation) int {
         return a.Bandwidth - b.Bandwidth
      },
   )
   for i, represent := range represents {
      if i >= 1 {
         fmt.Println()
      }
      fmt.Println(&represent)
   }
   return nil
}

const (
   widevine_system_id = "edef8ba979d64acea3c827dcd51d21ed"
   widevine_urn       = "urn:uuid:edef8ba9-79d6-4ace-a3c8-27dcd51d21ed"
)

func create(name string) (*os.File, error) {
   log.Println("Create", name)
   return os.Create(name)
}

func get(u *url.URL, head http.Header) ([]byte, error) {
   req := http.Request{Method: "GET", URL: u}
   if head != nil {
      req.Header = head
   } else {
      req.Header = http.Header{}
   }
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   switch resp.StatusCode {
   case http.StatusOK, http.StatusPartialContent:
   default:
      var data strings.Builder
      resp.Write(&data)
      return nil, errors.New(data.String())
   }
   return io.ReadAll(resp.Body)
}

func marshal(resp *http.Response) ([]byte, error) {
   var buf bytes.Buffer
   _, err := fmt.Fprintln(&buf, resp.Request.URL)
   if err != nil {
      return nil, err
   }
   err = resp.Write(&buf)
   if err != nil {
      return nil, err
   }
   return buf.Bytes(), nil
}

func unmarshal(data []byte) (*http.Response, error) {
   data1, data, _ := bytes.Cut(data, []byte{'\n'})
   var base url.URL
   err := base.UnmarshalBinary(data1)
   if err != nil {
      return nil, err
   }
   return http.ReadResponse(
      bufio.NewReader(bytes.NewReader(data)), &http.Request{URL: &base},
   )
}

func dash_create(represent *dash.Representation) (*os.File, error) {
   switch *represent.MimeType {
   case "audio/mp4":
      return create(".m4a")
   case "text/vtt":
      return create(".vtt")
   case "video/mp4":
      return create(".m4v")
   }
   return nil, errors.New(*represent.MimeType)
}

func init() {
   log.SetFlags(log.Ltime)
   http.DefaultClient.Transport = &transport{
      // x/net/http2: make Transport return nicer error when Amazon ALB hangs up
      // mid-response?
      // github.com/golang/go/issues/18639
      // x/net/http2: Transport ignores net/http.Transport.Proxy once connected
      // github.com/golang/go/issues/25793
      Protocols: &http.Protocols{},
      Proxy:     http.ProxyFromEnvironment,
   }
}

func (m *media_file) New(represent *dash.Representation) error {
   for _, content := range represent.ContentProtection {
      if content.SchemeIdUri == widevine_urn {
         if content.Pssh != "" {
            data, err := base64.StdEncoding.DecodeString(content.Pssh)
            if err != nil {
               return err
            }
            var box pssh.Box
            n, err := box.BoxHeader.Decode(data)
            if err != nil {
               return err
            }
            err = box.Read(data[n:])
            if err != nil {
               return err
            }
            m.pssh = box.Data
            break
         }
      }
   }
   return nil
}
type media_file struct {
   key_id    []byte // tenc
   pssh      []byte // pssh
   timescale uint64 // mdhd
   size      uint64 // trun
   duration  uint64 // trun
}

func (m *media_file) initialization(data []byte) ([]byte, error) {
   var file1 file.File
   err := file1.Read(data)
   if err != nil {
      return nil, err
   }
   // Moov
   moov, ok := file1.GetMoov()
   if !ok {
      return data, nil
   }
   // Moov.Pssh
   for _, pssh1 := range moov.Pssh {
      if pssh1.SystemId.String() == widevine_system_id {
         m.pssh = pssh1.Data
      }
      copy(pssh1.BoxHeader.Type[:], "free") // Firefox
   }
   // Moov.Trak
   m.timescale = uint64(moov.Trak.Mdia.Mdhd.Timescale)
   // Sinf
   sinf, ok := moov.Trak.Mdia.Minf.Stbl.Stsd.Sinf()
   if !ok {
      return data, nil
   }
   // Sinf.BoxHeader
   copy(sinf.BoxHeader.Type[:], "free") // Firefox
   // Sinf.Schi
   m.key_id = sinf.Schi.Tenc.DefaultKid[:]
   // SampleEntry
   sample, ok := moov.Trak.Mdia.Minf.Stbl.Stsd.SampleEntry()
   if !ok {
      return data, nil
   }
   // SampleEntry.BoxHeader
   copy(sample.BoxHeader.Type[:], sinf.Frma.DataFormat[:]) // Firefox
   return file1.Append(nil)
}
// segment can be VTT or anything
func (m *media_file) write_segment(data, key []byte) ([]byte, error) {
   if key == nil {
      return data, nil
   }
   var file1 file.File
   err := file1.Read(data)
   if err != nil {
      return nil, err
   }
   if m.duration/m.timescale < 10*60 {
      for _, sample := range file1.Moof.Traf.Trun.Sample {
         if sample.Duration == 0 {
            sample.Duration = file1.Moof.Traf.Tfhd.DefaultSampleDuration
         }
         m.duration += uint64(sample.Duration)
         if sample.Size == 0 {
            sample.Size = file1.Moof.Traf.Tfhd.DefaultSampleSize
         }
         m.size += uint64(sample.Size)
      }
      log.Println("bandwidth", m.timescale*m.size*8/m.duration)
   }
   if file1.Moof.Traf.Senc == nil {
      return data, nil
   }
   for i, data := range file1.Mdat.Data(&file1.Moof.Traf) {
      err = file1.Moof.Traf.Senc.Sample[i].Decrypt(data, key)
      if err != nil {
         return nil, err
      }
   }
   return file1.Append(nil)
}

func (p *progress) next() {
   p.segmentA++
   p.segmentB--
   timeB := time.Now().Unix()
   if timeB > p.timeB {
      log.Println(
         p.segmentB, "segment",
         p.durationB().Truncate(time.Second),
         "left",
      )
      p.timeB = timeB
   }
}

func (p *progress) set(segmentB int) {
   p.segmentB = segmentB
   p.timeA = time.Now()
   p.timeB = time.Now().Unix()
}

type progress struct {
   segmentA int64
   segmentB int
   timeA    time.Time
   timeB    int64
}

func (p *progress) durationA() time.Duration {
   return time.Since(p.timeA)
}

// keep last two terms separate
func (p *progress) durationB() time.Duration {
   return p.durationA() * time.Duration(p.segmentB) / time.Duration(p.segmentA)
}

type transport http.Transport

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
   if req.Header.Get("silent") == "" {
      log.Println(req.Method, req.URL)
   }
   return (*http.Transport)(t).RoundTrip(req)
}

func variation(value *dash.Representation, expect int) int {
   variation := value.Bandwidth - expect
   if variation < 0 {
      return -variation
   }
   return variation
}

func expected(values []dash.Representation, expect int) dash.Representation {
   a := values[0]
   for _, b := range values[1:] {
      if variation(&b, expect) < variation(&a, expect) {
         a = b
      }
   }
   return a
}

// github.com/golang/go/blob/go1.24.3/src/math/all_test.go#L2146
func tolerance(value *dash.Representation, expect int, percent float64) bool {
   return float64(variation(value, expect)) <= float64(expect)*percent
}
