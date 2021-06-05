package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
)

type Product struct {
	Name  string            `json:"name"`
	Img   string            `json:"img"`
	Attrs map[string]string `json:"attrs"`
}

func main() {
	c := colly.NewCollector(
		colly.Async(true),
	)

	extensions.RandomUserAgent(c)

	products := make([]Product, 0)

	c.OnHTML("ul.b-list__items_nofooter", func(element *colly.HTMLElement) { onHTML(&products, element) })
	c.OnResponse(onResponse)
	c.OnRequest(onRequest)

	query := "Niue Boba Fett 1oz"

	for i := 1; i <= 1; i++ {
		uri := url.URL{
			Scheme: "https",
			Host:   "www.ebay.com",
			Path:   "/sch/i.html",
		}

		q := uri.Query()
		q.Add("_from", "R40")
		q.Add("_nkw", query)
		q.Add("_sacat", "0") // Category ID
		q.Add("LH_PrefLoc", "1")
		q.Add("LH_Sold", "1")
		q.Add("LH_Complete", "1")
		q.Add("_udlo", "1")   // Low Price Range
		q.Add("_udhi", "100") // High Price Range
		q.Add("rt", "nc")
		q.Add("_ipg", "200")
		q.Add("_pgn", fmt.Sprintf("%d", i))

		uri.RawQuery = q.Encode()
		u := uri.String()

		err := c.Visit(u)
		if err != nil {
			log.Fatal(err)
			break
		}

		t := int(20*uniform()) * int(time.Second)
		time.Sleep(time.Duration(t))
	}

	c.Wait()

	file, err := os.OpenFile("result.json", os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(products)
	if err != nil {
		log.Fatal(err)
	}
}

func onHTML(products *[]Product, element *colly.HTMLElement) {
	if products == nil {
		return
	}

	element.ForEach("li.s-item", func(_ int, elem *colly.HTMLElement) {
		p := Product{Attrs: make(map[string]string)}
		p.Name = elem.ChildText("h3.s-item__title")

		img1 := elem.ChildAttr("img.s-item__image-img", "src")
		img2 := elem.ChildAttr("img.s-item__image-img", "data-src")

		if strings.Contains(img1, "ir.ebaystatic.com") {
			p.Img = img2
		} else {
			p.Img = img1
		}

		elem.ForEach("span.s-item__detail.s-item__detail--secondary", func(_ int, e *colly.HTMLElement) {
			attr := e.ChildText("span.s-item__dynamic")
			if attr == "" {
				return
			}

			kv := strings.SplitN(attr, ": ", 2)
			key := kv[0]
			val := kv[1]
			p.Attrs[key] = val

		})
		fmt.Printf("%+v", p)
		*products = append(*products, p)
	})
}

func onRequest(r *colly.Request) {
	log.Println("Visiting ", r.URL.String())

	uri := url.URL{
		Scheme: "https",
		Host:   "www.ebay.com",
		Path:   "/splashui/captcha",
	}
	q := uri.Query()
	q.Add("ap", "1")
	q.Add("appName", "orch")
	q.Add("ru", r.URL.String())
	fmt.Println(q.Encode())
	uri.RawQuery = q.Encode()

	r.Headers.Add("Cookie", "__uzma=153d2ce7-2b56-4b07-a0f9-67dbd700f16a; __uzmb=1609641196; __ssds=2; __ssuzjsr2=a9be0cd8e; __uzmaj2=00a67e99-3384-4391-989d-3bd46024cd6f; __uzmbj2=1609641197; cid=RJQzkFk3wdXkM8EF%23888252336; __uzme=3519; QuantumMetricUserID=2bd37bb459955dee3cf4da60f27929d1; D_HID=756BB114-57A3-3061-9DF8-BDE71FE4579A; D_IID=D8F94D47-FBE9-34C8-9D37-BE7AA4124D07; D_ZID=BEDBB469-A805-39DD-87F7-634DFF9A3194; D_SID=136.30.186.46; D_UID=4BFB31BA-2532-3460-A20F-1A9104230D68; D_ZUID=7B7F1DE0-B642-302C-BA02-30B06D47C299; ak_bmsc=05B873EEA3EC2E0FCC5AD979ED2FF8BA173724741125000055D65A60606AB377~plY8+D6y5MoM3R8P4k3rIkHhgZTwB6DHjmKNA9RfRzXOH8JVrgfg7lGMpuQMASL/Mv1Zf2hGeH6ArNBGH4xWyMxvW/LCzO7Nkgizv7uAzBMJ5KspliO0yms13tfM1mOtYHVBo+1y6bzfUUApb8pm077ZAmelktHJDzdm/Zc7ZjxFK/Il34xrfZ2Hh8Vth3sRHQ3760j8sSW2P4i4WrrM6PeFvaT0c+2xExwImGm3Ox+HY; bm_sv=7F49B02EB16C4F246E8CBC24E0F2A4B5~BrK4enLXyR9vj2/pMRhgqvGEo+mjVfeuWQn03eP3BbQSb2Nev9qGXiRZyP91KLdPDN2NCWEUc2oPPeqOCrtAbncuuwsYUpK6VYcUbP4JnK51w9BbgotIpvhkKSyDPQxLYsWfbC19cDELFowJlwCJYfR9d7V3mC4jqlmQlgiWpeM; __gads=ID=22a669b96385a75d:T=1616567106:S=ALNI_MadiVhftMTIvAvcXovRsBu8G4-8CQ; DG_SID=136.30.186.46:utVROTmCHo4YRg89JHs+rC7NlJwWgNCJuoTJJuDXh70; AMCVS_A71B5B5B54F607AB0A4C98A2%40AdobeOrg=1; DG_ZID=56B436DA-6F29-3177-87DB-B52CF327B536; DG_ZUID=27D8E81D-04A5-3F2F-9036-2307F98DEE0B; DG_HID=B14F6E97-3A4C-36C6-AEA2-99DCEB4B2D81; DG_IID=2D00D6C1-2C26-3F20-8DE9-4036412BD4E5; DG_UID=AE0E6545-E25C-3256-B9F1-CAA41368A6DD; ds1=ats/1616801548863; cssg=dfbb17c41760aad90773076cffad71d1; JSESSIONID=5A868D4C5BB649455FD5596A25131490; shs=BAQAAAXiVB3iwAAaAAVUAD2I/oowyMDQxMzMxNTMyMDA1LDI3EqX+YbI9WwO9P6xTQw6Mc3wkhg**; AMCV_A71B5B5B54F607AB0A4C98A2%40AdobeOrg=-408604571%7CMCMID%7C18721502967304683392404039235057215343%7CMCAAMLH-1618020143%7C7%7CMCAAMB-1618020143%7CRKhpRz8krg2tLO6pguXWp5olkAcUniQYPHaMWWgdJ3xzPWQmdj0y%7CMCCIDH%7C-682630538%7CMCOPTOUT-1617422543s%7CNONE%7CvVersion%7C4.6.0; __uzmc=5240174832516; __uzmd=1617417269; ds2=; __uzmcj2=8245639490559; __uzmdj2=1617417270; npii=btguid/dfbb17c41760aad90773076cffad71d1642a3b37^cguid/fd53d9ec1760aadc3ca082ddfbee5a4a642a3b37^; ebay=%5EsfLMD%3D0%5Esin%3Din%5Edv%3D6067ce56%5Esbf%3D%2320000000000010008000044%5Ecos%3D2%5Elrtjs%3D0.1%5Ecv%3D15555%5Ejs%3D1%5E; cpt=%5Ecpt_prvd%3Drecaptcha_v2%5Ecpt_guid%3D831b42ed-cf49-4217-9582-7d8535bbbbad%5E; ns1=BAQAAAXiVB3iwAAaAAKUADWJJCgwxODU3MzA5MzczLzA7ANgASmJJCgxjNjl8NjAxXjE2MTY4MDMzMjMzMDZeXjFeM3wyfDV8NHw3fDExXl5eNF4zXjEyXjEyXjJeMV4xXjBeMV4wXjFeNjQ0MjQ1OTA3NXVwsOIz73sZeZqvmhqshkxzOLjP; dp1=bkms/in642a3d8c^u1f/Colin642a3d8c^tzo/12c6067e49c^exc/0%3A0%3A1%3A1608f638c^pcid/88825233662490a0c^mpc/0%7C06075058c^u1p/ZWNvbW1ldF9sbGM*642a3d8c^bl/USen-US642a3d8c^expt/0001616801549533614f08cd^pbf/%232030040a000200050819c0200000462490a0c^; s=BAQAAAXiVB3iwAAWAAPgAIGBpKAxkZmJiMTdjNDE3NjBhYWQ5MDc3MzA3NmNmZmFkNzFkMQFFAAhiSQoMNjA1MWMwNzKSZPoCNcjaAWm9I3KSd4dQd7idKQ**; nonsession=BAQAAAXiVB3iwAAaAAJ0ACGJJCgwwMDAwMDAwMQFkAAdkKj2MIzAwMDAwYQAIABxgj2OMMTYxNzMyMjEwNHgyMzM5MTI5MjgwMTN4MHgyWQAzAA5iSQoMNjAwOTMtMTYyOSxVU0EAywACYGfdlDkyAEAAC2JJCgxlY29tbWV0X2xsYwAQAAtiSQoMZWNvbW1ldF9sbGMAygAgZCo9jGRmYmIxN2M0MTc2MGFhZDkwNzczMDc2Y2ZmYWQ3MWQxAAQAC2I/ooxlY29tbWV0X2xsYwCcADhiSQoMblkrc0haMlByQm1kajZ3Vm5ZK3NFWjJQckEyZGo2QU1sWVdpQ0ptQnBnK2RqNng5blkrc2VRPT3VPd9ALEJvAHeyfBDAJZXsOGrzvQ**")
	r.Headers.Add("Cache-Control", "max-age=0")
	r.Headers.Add("Accept-Language", "en-US,en;q=0.9")
	r.Headers.Add("Accept-Encoding", "gzip, deflate, br")
	r.Headers.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	r.Headers.Set("Referer", uri.String())
}

func onResponse(r *colly.Response) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewBuffer(r.Body))
	if err != nil {
		log.Fatal(err)
	}

	el := doc.Find("h3.s-item__title")
	if len(el.Nodes) == 0 {
		log.Println("No items found on ", r.Request.URL.String())
	}
}

func uniform() float64 {
	sig := rand.Uint64() % (1 << 52)
	return (1 + float64(sig)/(1<<52)) / math.Pow(2, geometric())
}

// geometric returns a number picked from a geometric
// distribution of parameter 0.5.
func geometric() float64 {
	b := 1
	for rand.Uint64()%2 == 0 {
		b++
	}
	return float64(b)
}
