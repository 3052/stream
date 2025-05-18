package net

import (
   "41.neocities.org/dash"
   "41.neocities.org/sofia/file"
   "41.neocities.org/sofia/pssh"
   "41.neocities.org/widevine"
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
   "path"
   "slices"
   "strings"
   "time"
)

type representation [1]dash.Representation

func (a *representation) expect(actual []representation, expect int) {
   *a = actual[0]
   for _, b := range actual[1:] {
      if b.variation(expect) < a.variation(expect) {
         *a = b
      }
   }
}

func (r *representation) variation(expect int) int {
   variation := r[0].Bandwidth - expect
   if variation < 0 {
      return -variation
   }
   return variation
}

// github.com/golang/go/blob/go1.24.3/src/math/all_test.go#L2146
func (r *representation) tolerance(expect int, percent float64) bool {
   return float64(r.variation(expect)) <= float64(expect)*percent
}

func (r Response) marshal() ([]byte, error) {
   var buf bytes.Buffer
   _, err := fmt.Fprintln(&buf, r[0].Request.URL)
   if err != nil {
      return nil, err
   }
   err = r[0].Write(&buf)
   if err != nil {
      return nil, err
   }
   return buf.Bytes(), nil
}

type Response [1]*http.Response

func (r *Response) unmarshal(data []byte) error {
   before, data, _ := bytes.Cut(data, []byte{'\n'})
   var base url.URL
   err := base.UnmarshalBinary(before)
   if err != nil {
      return err
   }
   r[0], err = http.ReadResponse(
      bufio.NewReader(bytes.NewReader(data)), &http.Request{URL: &base},
   )
   if err != nil {
      return err
   }
   return nil
}
func (r Response) Mpd(name string) error {
   data, err := r.marshal()
   if err != nil {
      return err
   }
   err = write_file(name, data)
   if err != nil {
      return err
   }
   err = r.unmarshal(data)
   if err != nil {
      return err
   }
   defer r[0].Body.Close()
   data, err = io.ReadAll(r[0].Body)
   if err != nil {
      return err
   }
   var mpd dash.Mpd
   err = mpd.Unmarshal(data)
   if err != nil {
      return err
   }
   represents := slices.SortedFunc(mpd.Representation(),
      func(a, b *dash.Representation) int {
         return a.Bandwidth - b.Bandwidth
      },
   )
   for i, represent := range represents {
      if i >= 1 {
         fmt.Println()
      }
      fmt.Println(represent)
   }
   return nil
}

func (e *License) Download(name, id string) error {
   data, err := os.ReadFile(name)
   if err != nil {
      return err
   }
   var resp Response
   err = resp.unmarshal(data)
   if err != nil {
      return err
   }
   defer resp[0].Body.Close()
   data, err = io.ReadAll(resp[0].Body)
   if err != nil {
      return err
   }
   var mpd dash.Mpd
   err = mpd.Unmarshal(data)
   if err != nil {
      return err
   }
   mpd.Set(resp[0].Request.URL)
   for represent := range mpd.Representation() {
      if represent.Id == id {
         if represent.SegmentBase != nil {
            return e.segment_base(represent)
         }
         if represent.SegmentList != nil {
            return e.segment_list(represent)
         }
         return e.segment_template(represent)
      }
   }
   return nil
}

func write_file(name string, data []byte) error {
   err := os.MkdirAll(path.Dir(name), os.ModePerm)
   if err != nil {
      return err
   }
   log.Println("WriteFile", name)
   return os.WriteFile(name, data, os.ModePerm)
}

var ThreadCount = 1

const (
   widevine_system_id = "edef8ba979d64acea3c827dcd51d21ed"
   widevine_urn       = "urn:uuid:edef8ba9-79d6-4ace-a3c8-27dcd51d21ed"
)

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

func create(name string) (*os.File, error) {
   log.Println("Create", name)
   return os.Create(name)
}

type License struct {
   ClientId   string
   PrivateKey string
   Widevine   func([]byte) ([]byte, error)
}

func (e *License) segment_template(represent *dash.Representation) error {
   var media media_file
   err := media.New(represent)
   if err != nil {
      return err
   }
   file1, err := dash_create(represent)
   if err != nil {
      return err
   }
   defer file1.Close()
   if initial := represent.SegmentTemplate.Initialization; initial != "" {
      address, err := initial.Url(represent)
      if err != nil {
         return err
      }
      data1, err := get_segment(address, nil)
      if err != nil {
         return err
      }
      data1, err = media.initialization(data1)
      if err != nil {
         return err
      }
      _, err = file1.Write(data1)
      if err != nil {
         return err
      }
   }
   key, err := e.get_key(&media)
   if err != nil {
      return err
   }
   var segments []int
   for represent1 := range represent.Representation() {
      segments = slices.AppendSeq(segments, represent1.Segment())
   }
   var progress1 progress
   progress1.set(len(segments))
   for chunk := range slices.Chunk(segments, ThreadCount) {
      var (
         datas = make([][]byte, len(chunk))
         errs = make(chan error)
      )
      for i, segment := range chunk {
         address, err := represent.SegmentTemplate.Media.Url(represent, segment)
         if err != nil {
            return err
         }
         go func() {
            datas[i], err = get_segment(address, nil)
            errs <- err
            progress1.next()
         }()
      }
      for range chunk {
         err := <-errs
         if err != nil {
            return err
         }
      }
      for _, data := range datas {
         data, err = media.write_segment(data, key)
         if err != nil {
            return err
         }
         _, err = file1.Write(data)
         if err != nil {
            return err
         }
      }
   }
   return nil
}

func (e *License) segment_base(represent *dash.Representation) error {
   if ThreadCount != 1 {
      return errors.New("ThreadCount")
   }
   var media media_file
   err := media.New(represent)
   if err != nil {
      return err
   }
   file1, err := dash_create(represent)
   if err != nil {
      return err
   }
   defer file1.Close()
   data, err := get_segment(represent.BaseUrl[0], http.Header{
      "range": {"bytes=" + represent.SegmentBase.Initialization.Range},
   })
   if err != nil {
      return err
   }
   data, err = media.initialization(data)
   if err != nil {
      return err
   }
   _, err = file1.Write(data)
   if err != nil {
      return err
   }
   key, err := e.get_key(&media)
   if err != nil {
      return err
   }
   data, err = get_segment(represent.BaseUrl[0], http.Header{
      "range": {"bytes=" + represent.SegmentBase.IndexRange},
   })
   if err != nil {
      return err
   }
   var file2 file.File
   err = file2.Read(data)
   if err != nil {
      return err
   }
   var progress1 progress
   progress1.set(len(file2.Sidx.Reference))
   head := http.Header{}
   var index dash.Range
   err = index.Set(represent.SegmentBase.IndexRange)
   if err != nil {
      return err
   }
   for _, reference := range file2.Sidx.Reference {
      index[0] = index[1] + 1
      index[1] += uint64(reference.Size())
      head.Set("range", "bytes="+ index.String())
      data, err = get_segment(represent.BaseUrl[0], head)
      if err != nil {
         return err
      }
      progress1.next()
      data, err = media.write_segment(data, key)
      if err != nil {
         return err
      }
      _, err = file1.Write(data)
      if err != nil {
         return err
      }
   }
   return nil
}

func (e *License) segment_list(represent *dash.Representation) error {
   if ThreadCount != 1 {
      return errors.New("ThreadCount")
   }
   var media media_file
   err := media.New(represent)
   if err != nil {
      return err
   }
   file1, err := dash_create(represent)
   if err != nil {
      return err
   }
   defer file1.Close()
   data, err := get_segment(
      represent.SegmentList.Initialization.SourceUrl[0], nil,
   )
   if err != nil {
      return err
   }
   data, err = media.initialization(data)
   if err != nil {
      return err
   }
   _, err = file1.Write(data)
   if err != nil {
      return err
   }
   key, err := e.get_key(&media)
   if err != nil {
      return err
   }
   var progress1 progress
   progress1.set(len(represent.SegmentList.SegmentUrl))
   for _, segment := range represent.SegmentList.SegmentUrl {
      data, err := get_segment(segment.Media[0], nil)
      if err != nil {
         return err
      }
      progress1.next()
      data, err = media.write_segment(data, key)
      if err != nil {
         return err
      }
      _, err = file1.Write(data)
      if err != nil {
         return err
      }
   }
   return nil
}

func init() {
   log.SetFlags(log.Ltime)
   http.DefaultClient.Transport = transport{}
}

// LOG
// PROXY
type transport struct{}

func (transport) RoundTrip(req *http.Request) (*http.Response, error) {
   log.Println(req.Method, req.URL)
   return http.DefaultTransport.RoundTrip(req)
}

// NO LOG
// NO PROXY
var Client = http.Client{
   Transport: &http.Transport{},
}

func get_segment(u *url.URL, head http.Header) ([]byte, error) {
   req := http.Request{Method: "GET", URL: u}
   if head != nil {
      req.Header = head
   } else {
      req.Header = http.Header{}
   }
   resp, err := Client.Do(&req)
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

func (e *License) get_key(media *media_file) ([]byte, error) {
   if media.key_id == nil {
      return nil, nil
   }
   private_key, err := os.ReadFile(e.PrivateKey)
   if err != nil {
      return nil, err
   }
   client_id, err := os.ReadFile(e.ClientId)
   if err != nil {
      return nil, err
   }
   if media.pssh == nil {
      var pssh1 widevine.Pssh
      pssh1.KeyIds = [][]byte{media.key_id}
      media.pssh = pssh1.Marshal()
   }
   log.Println("PSSH", base64.StdEncoding.EncodeToString(media.pssh))
   var module widevine.Cdm
   err = module.New(private_key, client_id, media.pssh)
   if err != nil {
      return nil, err
   }
   data, err := module.RequestBody()
   if err != nil {
      return nil, err
   }
   data, err = e.Widevine(data)
   if err != nil {
      return nil, err
   }
   var body widevine.ResponseBody
   err = body.Unmarshal(data)
   if err != nil {
      return nil, err
   }
   block, err := module.Block(body)
   if err != nil {
      return nil, err
   }
   for container := range body.Container() {
      if bytes.Equal(container.Id(), media.key_id) {
         key := container.Key(block)
         log.Println("key", base64.StdEncoding.EncodeToString(key))
         var zero [16]byte
         if !bytes.Equal(key, zero[:]) {
            return key, nil
         }
      }
   }
   return nil, errors.New("get_key")
}
