package controllers

import (
	"fmt"
	"net/http"
)

type ImageWS struct{}

func (handler *ImageWS) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	//bts, _ := ioutil.ReadAll(r.Body)
	//fmt.Println(string(bts))
	//fmt.Println("111")

	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Content-Type", "text/xml; charset=utf-8")
	//a := []byte("<?xml version=\"1.0\" encoding=\"utf-8\"?><soap:Envelope xmlns:soap=\"http://schemas.xmlsoap.org/soap/envelope/\" xmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\" xmlns:xsd=\"http://www.w3.org/2001/XMLSchema\"><soap:Body><RequestResponse xmlns=\"http://tempuri.org/\"><RequestResult>&lt;?xml version=\"1.0\"?&gt;&lt;Response&gt;  &lt;Header ServerIP=\"10.2.6.250\" ShouldRecordPerformanceTime=\"false\" UserID=\"900407\" RequestID=\"f2450c88-85d7-4b82-b2b4-12f4379bd0b3\" ResultCode=\"Success\" AssemblyVersion=\"1..2.6\" RequestBodySize=\"0\" SerializeMode=\"Xml\" RouteStep=\"1\" Environment=\"fws\" /&gt;  &lt;SaveResponse&gt;    &lt;CheckPass&gt;true&lt;/CheckPass&gt;    &lt;OriginalPath&gt;\\fd\\tuangou\\group1\\M01\\6F\\08\\CgIG6VXfwiaAWJfmAANHCVkEN94713.jpg&lt;/OriginalPath&gt;    &lt;TargetPath&gt;\\fd\\tuangou\\group1\\M01\\6F\\08\\CgIG6VXfwiaAWJfmAANHCVkEN94713.jpg&lt;/TargetPath&gt;    &lt;Process&gt;      &lt;ProcessResponse&gt;        &lt;ID&gt;1&lt;/ID&gt;        &lt;Path&gt;\\fd\\tuangou\\group1\\M01\\6F\\08\\CgIG6VXfwiaAWJfmAANHCVkEN94713_20_20.jpg&lt;/Path&gt;        &lt;Info&gt;20,20&lt;/Info&gt;     &lt;/ProcessResponse&gt;    &lt;/Process&gt;  &lt;/SaveResponse&gt; &lt;/Response&gt;</RequestResult></RequestResponse></soap:Body></soap:Envelope>")
	content := []byte("")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(content)))
	w.Write(content)
}
