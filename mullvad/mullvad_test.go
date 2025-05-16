package mullvad

import (
   "log"
   "net/http"
   "testing"
)

// ca-yyc-wg-201
const paramount = "https://www.paramountplus.com/apps-api/v2.0/androidphone/video/cid/Y8sKvb2bIoeX4XZbsfjadF4GhNPwcjTQ.json?at=ABAAAAAAAAAAAAAAAAAAAAAAU%2FEmq2DGAxNGYi71fZdi2mb6UTU2%2BelcD3bGEhs6bo4%3D"

func Test(t *testing.T) {
   http.DefaultClient.Transport = &Transport{}
   log.SetFlags(log.Ltime)
   defer Disconnect()
   req, err := http.NewRequest("HEAD", paramount, nil)
   if err != nil {
      t.Fatal(err)
   }
   req.Header.Set("proxy", "true")
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      t.Fatal(err)
   }
   if resp.StatusCode != http.StatusOK {
      t.Fatal(resp.Status)
   }
}
