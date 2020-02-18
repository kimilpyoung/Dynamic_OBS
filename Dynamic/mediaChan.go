package Dynamic
type commChan struct {
  cmd string
  body string
  sdp_type string
  sdp_des string
  media *commMedia
}

type commMedia struct {
  audio bool
  data *[]byte
  ts uint64
  duration float64
}
