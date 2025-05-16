package net

import (
   "41.neocities.org/dash"
   "41.neocities.org/sofia/file"
   "41.neocities.org/widevine"
   "bytes"
   "encoding/base64"
   "errors"
   "io"
   "log"
   "net/http"
   "os"
   "slices"
)

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
      data1, err := get(address, nil)
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
   head := http.Header{}
   head.Set("silent", "true")
   var segments []int
   for r := range represent.Representation() {
      segments = slices.AppendSeq(segments, r.Segment())
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
            datas[i], err = get(address, head)
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
   data, err := get(represent.BaseUrl[0], http.Header{
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
   data, err = get(represent.BaseUrl[0], http.Header{
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
   head.Set("silent", "true")
   var index dash.Range
   err = index.Set(represent.SegmentBase.IndexRange)
   if err != nil {
      return err
   }
   for _, reference := range file2.Sidx.Reference {
      index[0] = index[1] + 1
      index[1] += uint64(reference.Size())
      head.Set("range", "bytes="+ index.String())
      data, err = get(represent.BaseUrl[0], head)
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
   data, err := get(represent.SegmentList.Initialization.SourceUrl[0], nil)
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
   head := http.Header{}
   head.Set("silent", "true")
   for _, segment := range represent.SegmentList.SegmentUrl {
      data, err := get(segment.Media[0], head)
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

type License struct {
   ClientId   string
   PrivateKey string
   Widevine   func([]byte) ([]byte, error)
}

func (e *License) Download(name, id string) error {
   data, err := os.ReadFile(name)
   if err != nil {
      return err
   }
   resp, err := unmarshal(data)
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
   mpd1.Set(resp.Request.URL)
   for represent := range mpd1.Representation() {
      if represent.Id == id {
         if represent.SegmentBase != nil {
            return e.segment_base(&represent)
         }
         if represent.SegmentList != nil {
            return e.segment_list(&represent)
         }
         return e.segment_template(&represent)
      }
   }
   return nil
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
