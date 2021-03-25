package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pborman/uuid"
	kafka "github.com/segmentio/kafka-go"
)

type User struct {
	Conn  *websocket.Conn
	ID    string
	ticks []tick
}

type tick struct {
	name  string
	value string
}

var users = make(map[string]*User)
var mutex sync.RWMutex
var ticks_map = make(map[string]string)

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Home Page")
}

func getKafkaReader(kafkaURL,
	topic, groupID string) *kafka.Reader {
	brokers := strings.Split(kafkaURL, ",")
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:     brokers,
		GroupID:     groupID,
		Topic:       topic,
		StartOffset: -1,
	})
}

func fullyTicksMap() {
	for i := 0; i < len(Ticks); i++ {
		ticks_map[Ticks[i].name] = ""
	}
}
func updateTick(values []string,
	ti []tick) {
	for i := 0; i < len(values); i++ {
		if values[i] != ti[i].value {
			ti[i].value = values[i]
			ticks_map[ti[i].name] = ti[i].value
		}
	}
}

func wsEndpoint(w http.ResponseWriter, r *http.Request, ch chan string, done chan bool) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	uid := uuid.New()
	user := new(User)
	user.ID = uid
	user.Conn = ws
	ticks_user := r.URL.Query()["ticks"]
	ticks_user_array := []tick{}
	for t := range ticks_user {
		t_i := new(tick)
		t_i.name = ticks_user[t]
		t_i.value = ""
		ticks_user_array = append(ticks_user_array, *t_i)
	}
	user.ticks = ticks_user_array

	mutex.Lock()
	users[uid] = user
	mutex.Unlock()
	fmt.Println(user)
	for {
		select {
		case <-done:
			return
		case _ = <-ch:
			for u := range users {
				for pos := range users[u].ticks {
					t_name_user := users[u].ticks[pos].name
					t_value_user := users[u].ticks[pos].value
					if t_value_user != ticks_map[t_name_user] {
						users[u].ticks[pos].value = ticks_map[t_name_user]

						err = users[u].Conn.WriteMessage(1, []byte(fmt.Sprintf("%s - %s - %d", t_name_user, ticks_map[t_name_user], (time.Now().UnixNano()/int64(time.Millisecond)-10800000))))
						if err != nil {
							fmt.Println(err)
						}
					}
				}
			}

		}

	}
}
func setupRoutes(ch chan string, done chan bool) {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		wsEndpoint(w, r, ch, done)
	})
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {
	fullyTicksMap()
	fmt.Println("Hello World")
	flagHost := flag.String("host", "18.228.170.236:9092", "string")
	flagTopic := flag.String("topic", "b3-simulate", "string")
	flagGroup := flag.String("group", "b3-simulate-c1", "string")
	flag.Parse()

	kafkaURL := *flagHost
	topic := *flagTopic
	groupID := *flagGroup

	reader := getKafkaReader(kafkaURL,
		topic, groupID)

	defer reader.Close()
	ticker := time.NewTicker(300 * time.Millisecond)
	done := make(chan bool)
	values := make(chan string)
	go func(values chan string) {
		for {
			select {
			case <-done:
				return
			case _ = <-ticker.C:
				m, err := reader.ReadMessage(context.Background())
				if err != nil {
					log.Fatalln(err)
				}
				parsed := strings.Split(string(m.Value), ";")
				updateTick(parsed, Ticks)
				values <- "updated"
			}
		}
	}(values)
	setupRoutes(values, done)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

var Ticks = []tick{{name: "AALR3", value: ""}, {name: "AAPL34", value: ""}, {name: "ABBV34", value: ""}, {name: "ABCB4", value: ""}, {name: "ABEV3", value: ""}, {name: "ADHM3", value: ""}, {name: "AFLT3", value: ""}, {name: "AGRO3", value: ""}, {name: "ALPA3", value: ""}, {name: "ALPA4", value: ""}, {name: "ALSC3", value: ""}, {name: "ALUP3", value: ""}, {name: "ALUP4", value: ""}, {name: "ALUP11", value: ""}, {name: "AMAR3", value: ""}, {name: "AMGN34", value: ""}, {name: "AMZO34", value: ""}, {name: "ANIM3", value: ""}, {name: "APER3", value: ""}, {name: "ARZZ3", value: ""}, {name: "ATOM3", value: ""}, {name: "ATTB34", value: ""}, {name: "AZEV4", value: ""}, {name: "AZUL4", value: ""}, {name: "B3SA3", value: ""}, {name: "BAHI3", value: ""}, {name: "BAUH4", value: ""}, {name: "BAZA3", value: ""}, {name: "BBAS3", value: ""}, {name: "BBDC3", value: ""}, {name: "BBDC4", value: ""}, {name: "BBRK3", value: ""}, {name: "BBSE3", value: ""}, {name: "BEEF3", value: ""}, {name: "BEES3", value: ""}, {name: "BEES4", value: ""}, {name: "BERK34", value: ""}, {name: "BGIP4", value: ""}, {name: "BIDI4", value: ""}, {name: "BIIB34", value: ""}, {name: "BIOM3", value: ""}, {name: "BKBR3", value: ""}, {name: "BMEB4", value: ""}, {name: "BNBR3", value: ""}, {name: "BOAC34", value: ""}, {name: "BOBR4", value: ""}, {name: "BOEI34", value: ""}, {name: "BOXP34", value: ""}, {name: "BPAC3", value: ""}, {name: "BPAC5", value: ""}, {name: "BPAC11", value: ""}, {name: "BPAN4", value: ""}, {name: "BRAP3", value: ""}, {name: "BRAP4", value: ""}, {name: "BRDT3", value: ""}, {name: "BRFS3", value: ""}, {name: "BRIV4", value: ""}, {name: "BRKM3", value: ""}, {name: "BRKM5", value: ""}, {name: "BRML3", value: ""}, {name: "BRPR3", value: ""}, {name: "BRSR3", value: ""}, {name: "BRSR6", value: ""}, {name: "BSEV3", value: ""}, {name: "BTOW3", value: ""}, {name: "BTTL3", value: ""}, {name: "CAMB4", value: ""}, {name: "CAML3", value: ""}, {name: "CARD3", value: ""}, {name: "CBEE3", value: ""}, {name: "CCPR3", value: ""}, {name: "CCRO3", value: ""}, {name: "CCXC3", value: ""}, {name: "CEDO4", value: ""}, {name: "CEEB3", value: ""}, {name: "CELP3", value: ""}, {name: "CELP5", value: ""}, {name: "CELP7", value: ""}, {name: "CESP3", value: ""}, {name: "CESP5", value: ""}, {name: "CESP6", value: ""}, {name: "CGAS3", value: ""}, {name: "CGAS5", value: ""}, {name: "CGRA3", value: ""}, {name: "CGRA4", value: ""}, {name: "CHVX34", value: ""}, {name: "CIEL3", value: ""}, {name: "CLGN34", value: ""}, {name: "CLSC4", value: ""}, {name: "CMCS34", value: ""}, {name: "CMIG3", value: ""}, {name: "CMIG4", value: ""}, {name: "COCA34", value: ""}, {name: "COCE3", value: ""}, {name: "COCE5", value: ""}, {name: "COLG34", value: ""}, {name: "COPH34", value: ""}, {name: "COWC34", value: ""}, {name: "CPFE3", value: ""}, {name: "CPLE3", value: ""}, {name: "CPLE6", value: ""}, {name: "CPRE3", value: ""}, {name: "CRDE3", value: ""}, {name: "CRFB3", value: ""}, {name: "CRIV3", value: ""}, {name: "CRPG5", value: ""}, {name: "CRPG6", value: ""}, {name: "CSAN3", value: ""}, {name: "CSMG3", value: ""}, {name: "CSNA3", value: ""}, {name: "CSRN3", value: ""}, {name: "CSRN5", value: ""}, {name: "CTGP34", value: ""}, {name: "CTKA4", value: ""}, {name: "CTNM3", value: ""}, {name: "CTNM4", value: ""}, {name: "CTSA3", value: ""}, {name: "CTSA4", value: ""}, {name: "CTSH34", value: ""}, {name: "CVCB3", value: ""}, {name: "CVSH34", value: ""}, {name: "CYRE3", value: ""}, {name: "DEAI34", value: ""}, {name: "DIRR3", value: ""}, {name: "DISB34", value: ""}, {name: "DMMO3", value: ""}, {name: "DTEX3", value: ""}, {name: "EALT4", value: ""}, {name: "ECOR3", value: ""}, {name: "EGIE3", value: ""}, {name: "EKTR4", value: ""}, {name: "ELEK4", value: ""}, {name: "ELET3", value: ""}, {name: "ELET5", value: ""}, {name: "ELET6", value: ""}, {name: "ELPL3", value: ""}, {name: "EMAE4", value: ""}, {name: "EMBR3", value: ""}, {name: "ENBR3", value: ""}, {name: "ENEV3", value: ""}, {name: "ENGI3", value: ""}, {name: "ENGI4", value: ""}, {name: "ENGI11", value: ""}, {name: "ENMA3B", value: ""}, {name: "ENMT3", value: ""}, {name: "EQTL3", value: ""}, {name: "ESTC3", value: ""}, {name: "EUCA4", value: ""}, {name: "EVEN3", value: ""}, {name: "EXXO34", value: ""}, {name: "EZTC3", value: ""}, {name: "FBOK34", value: ""}, {name: "FESA3", value: ""}, {name: "FESA4", value: ""}, {name: "FHER3", value: ""}, {name: "FIBR3", value: ""}, {name: "FJTA3", value: ""}, {name: "FJTA4", value: ""}, {name: "FLRY3", value: ""}, {name: "FRAS3", value: ""}, {name: "GBIO33", value: ""}, {name: "GFSA3", value: ""}, {name: "GGBR3", value: ""}, {name: "GGBR4", value: ""}, {name: "GILD34", value: ""}, {name: "GNDI3", value: ""}, {name: "GOAU3", value: ""}, {name: "GOAU4", value: ""}, {name: "GOGL34", value: ""}, {name: "GOGL35", value: ""}, {name: "GOLL4", value: ""}, {name: "GPIV33", value: ""}, {name: "GPSI34", value: ""}, {name: "GRND3", value: ""}, {name: "GSGI34", value: ""}, {name: "GSHP3", value: ""}, {name: "GUAR3", value: ""}, {name: "GUAR4", value: ""}, {name: "HAGA4", value: ""}, {name: "HAPV3", value: ""}, {name: "HBOR3", value: ""}, {name: "HBTS5", value: ""}, {name: "HETA4", value: ""}, {name: "HGTX3", value: ""}, {name: "HOME34", value: ""}, {name: "HPQB34", value: ""}, {name: "HYPE3", value: ""}, {name: "IBMB34", value: ""}, {name: "IDNT3", value: ""}, {name: "IDVL3", value: ""}, {name: "IDVL4", value: ""}, {name: "IGTA3", value: ""}, {name: "IRBR3", value: ""}, {name: "ITLC34", value: ""}, {name: "ITSA3", value: ""}, {name: "ITSA4", value: ""}, {name: "ITUB3", value: ""}, {name: "ITUB4", value: ""}, {name: "JBDU3", value: ""}, {name: "JBDU4", value: ""}, {name: "JBSS3", value: ""}, {name: "JFEN3", value: ""}, {name: "JHSF3", value: ""}, {name: "JNJB34", value: ""}, {name: "JPSA3", value: ""}, {name: "JSLG3", value: ""}, {name: "KEPL3", value: ""}, {name: "KHCB34", value: ""}, {name: "KLBN3", value: ""}, {name: "KLBN4", value: ""}, {name: "KLBN11", value: ""}, {name: "KROT3", value: ""}, {name: "LAME3", value: ""}, {name: "LAME4", value: ""}, {name: "LCAM3", value: ""}, {name: "LEVE3", value: ""}, {name: "LIGT3", value: ""}, {name: "LILY34", value: ""}, {name: "LINX3", value: ""}, {name: "LIQO3", value: ""}, {name: "LLIS3", value: ""}, {name: "LMTB34", value: ""}, {name: "LOGG3", value: ""}, {name: "LOGN3", value: ""}, {name: "LPSB3", value: ""}, {name: "LREN3", value: ""}, {name: "MACY34", value: ""}, {name: "MAGG3", value: ""}, {name: "MCDC34", value: ""}, {name: "MDIA3", value: ""}, {name: "MDLZ34", value: ""}, {name: "MDTC34", value: ""}, {name: "MEAL3", value: ""}, {name: "METB34", value: ""}, {name: "MGEL4", value: ""}, {name: "MGLU3", value: ""}, {name: "MILS3", value: ""}, {name: "MMMC34", value: ""}, {name: "MNDL3", value: ""}, {name: "MNPR3", value: ""}, {name: "MOSC34", value: ""}, {name: "MOVI3", value: ""}, {name: "MPLU3", value: ""}, {name: "MRCK34", value: ""}, {name: "MRFG3", value: ""}, {name: "MRVE3", value: ""}, {name: "MSBR34", value: ""}, {name: "MSCD34", value: ""}, {name: "MSFT34", value: ""}, {name: "MTSA4", value: ""}, {name: "MULT3", value: ""}, {name: "MYPK3", value: ""}, {name: "NATU3", value: ""}, {name: "NFLX34", value: ""}, {name: "NIKE34", value: ""}, {name: "ODPV3", value: ""}, {name: "OFSA3", value: ""}, {name: "OGXP3", value: ""}, {name: "OMGE3", value: ""}, {name: "PARD3", value: ""}, {name: "PCAR4", value: ""}, {name: "PEAB3", value: ""}, {name: "PEPB34", value: ""}, {name: "PETR3", value: ""}, {name: "PETR4", value: ""}, {name: "PFIZ34", value: ""}, {name: "PFRM3", value: ""}, {name: "PGCO34", value: ""}, {name: "PINE4", value: ""}, {name: "PLAS3", value: ""}, {name: "PMAM3", value: ""}, {name: "PNVL3", value: ""}, {name: "PNVL4", value: ""}, {name: "POMO3", value: ""}, {name: "POMO4", value: ""}, {name: "POSI3", value: ""}, {name: "PPLA11", value: ""}, {name: "PRIO3", value: ""}, {name: "PSSA3", value: ""}, {name: "PTBL3", value: ""}, {name: "PTNT4", value: ""}, {name: "QGEP3", value: ""}, {name: "QUAL3", value: ""}, {name: "RADL3", value: ""}, {name: "RAIL3", value: ""}, {name: "RANI3", value: ""}, {name: "RANI4", value: ""}, {name: "RAPT3", value: ""}, {name: "RAPT4", value: ""}, {name: "RCSL3", value: ""}, {name: "RCSL4", value: ""}, {name: "RDNI3", value: ""}, {name: "REDE3", value: ""}, {name: "RENT3", value: ""}, {name: "RLOG3", value: ""}, {name: "RNEW3", value: ""}, {name: "RNEW4", value: ""}, {name: "RNEW11", value: ""}, {name: "ROMI3", value: ""}, {name: "ROST34", value: ""}, {name: "RSID3", value: ""}, {name: "SANB3", value: ""}, {name: "SANB4", value: ""}, {name: "SANB11", value: ""}, {name: "SAPR3", value: ""}, {name: "SAPR4", value: ""}, {name: "SAPR11", value: ""}, {name: "SBSP3", value: ""}, {name: "SBUB34", value: ""}, {name: "SCAR3", value: ""}, {name: "SCHW34", value: ""}, {name: "SEDU3", value: ""}, {name: "SEER3", value: ""}, {name: "SGPS3", value: ""}, {name: "SHOW3", value: ""}, {name: "SHUL4", value: ""}, {name: "SLCE3", value: ""}, {name: "SMLS3", value: ""}, {name: "SMTO3", value: ""}, {name: "SNSL3", value: ""}, {name: "SSBR3", value: ""}, {name: "STBP3", value: ""}, {name: "SULA11", value: ""}, {name: "SUZB3", value: ""}, {name: "TAEE3", value: ""}, {name: "TAEE4", value: ""}, {name: "TAEE11", value: ""}, {name: "TCSA3", value: ""}, {name: "TECN3", value: ""}, {name: "TELB3", value: ""}, {name: "TELB4", value: ""}, {name: "TEND3", value: ""}, {name: "TESA3", value: ""}, {name: "TEXA34", value: ""}, {name: "TGMA3", value: ""}, {name: "TIET3", value: ""}, {name: "TIET4", value: ""}, {name: "TIET11", value: ""}, {name: "TIMP3", value: ""}, {name: "TKNO4", value: ""}, {name: "TOTS3", value: ""}, {name: "TOYB4", value: ""}, {name: "TRIS3", value: ""}, {name: "TRPL3", value: ""}, {name: "TRPL4", value: ""}, {name: "TRPN3", value: ""}, {name: "TUPY3", value: ""}, {name: "TXRX4", value: ""}, {name: "UCAS3", value: ""}, {name: "UGPA3", value: ""}, {name: "UNIP3", value: ""}, {name: "UNIP5", value: ""}, {name: "UNIP6", value: ""}, {name: "UPAC34", value: ""}, {name: "UPSS34", value: ""}, {name: "USIM3", value: ""}, {name: "USIM5", value: ""}, {name: "USIM6", value: ""}, {name: "USSX34", value: ""}, {name: "VALE3", value: ""}, {name: "VERZ34", value: ""}, {name: "VISA34", value: ""}, {name: "VIVT3", value: ""}, {name: "VIVT4", value: ""}, {name: "VLID3", value: ""}, {name: "VLOE34", value: ""}, {name: "VULC3", value: ""}, {name: "VVAR3", value: ""}, {name: "WALM34", value: ""}, {name: "WEGE3", value: ""}, {name: "WHRL3", value: ""}, {name: "WIZS3", value: ""}, {name: "WLMM4", value: ""}, {name: "WSON33", value: ""}, {name: "WUNI34", value: ""}, {name: "XRXB34", value: ""}, {name: "TPIS3", value: ""}, {name: "BPHA3", value: ""}, {name: "ETER3", value: ""}, {name: "FRTA3", value: ""}, {name: "GPCP3", value: ""}, {name: "IGBR3", value: ""}, {name: "INEP3", value: ""}, {name: "INEP4", value: ""}, {name: "LUPA3", value: ""}, {name: "MMXM3", value: ""}, {name: "MWET4", value: ""}, {name: "OIBR3", value: ""}, {name: "OIBR4", value: ""}, {name: "OSXB3", value: ""}, {name: "PDGR3", value: ""}, {name: "SLED3", value: ""}, {name: "SLED4", value: ""}, {name: "TCNO3", value: ""}, {name: "TCNO4", value: ""}, {name: "TEKA4", value: ""}, {name: "VIVR3", value: ""}, {name: "PLAS1", value: ""}, {name: "POMO2", value: ""}, {name: "VIVR1", value: ""}, {name: "ABCP11", value: ""}, {name: "AEFI11", value: ""}, {name: "AGCX11", value: ""}, {name: "ALMI11", value: ""}, {name: "ALZR11", value: ""}, {name: "BBFI11B", value: ""}, {name: "BBPO11", value: ""}, {name: "BBRC11", value: ""}, {name: "BBVJ11", value: ""}, {name: "BCFF11", value: ""}, {name: "BCIA11", value: ""}, {name: "BCRI11", value: ""}, {name: "BMLC11B", value: ""}, {name: "BNFS11", value: ""}, {name: "BPFF11", value: ""}, {name: "BRCR11", value: ""}, {name: "BZLI11", value: ""}, {name: "CARE11", value: ""}, {name: "CBOP11", value: ""}, {name: "CEOC11", value: ""}, {name: "CNES11", value: ""}, {name: "CPTS11B", value: ""}, {name: "CTXT11", value: ""}, {name: "CXCE11B", value: ""}, {name: "CXRI11", value: ""}, {name: "CXTL11", value: ""}, {name: "DRIT11B", value: ""}, {name: "EDFO11B", value: ""}, {name: "EDGA11", value: ""}, {name: "EURO11", value: ""}, {name: "FAED11", value: ""}, {name: "FAMB11B", value: ""}, {name: "FCFL11", value: ""}, {name: "FEXC11", value: ""}, {name: "FFCI11", value: ""}, {name: "FIGS11", value: ""}, {name: "FIIB11", value: ""}, {name: "FIIP11B", value: ""}, {name: "FIVN11", value: ""}, {name: "FIXX11", value: ""}, {name: "FLMA11", value: ""}, {name: "FLRP11", value: ""}, {name: "FMOF11", value: ""}, {name: "FPAB11", value: ""}, {name: "FVBI11", value: ""}, {name: "FVPQ11", value: ""}, {name: "GGRC11", value: ""}, {name: "GRLV11", value: ""}, {name: "HCRI11", value: ""}, {name: "HFOF11", value: ""}, {name: "HGBS11", value: ""}, {name: "HGCR11", value: ""}, {name: "HGJH11", value: ""}, {name: "HGLG11", value: ""}, {name: "HGRE11", value: ""}, {name: "HGRU11", value: ""}, {name: "HMOC11", value: ""}, {name: "HTMX11", value: ""}, {name: "HUSC11", value: ""}, {name: "IRDM11", value: ""}, {name: "JRDM11", value: ""}, {name: "JSRE11", value: ""}, {name: "KNCR11", value: ""}, {name: "KNHY11", value: ""}, {name: "KNIP11", value: ""}, {name: "KNRE11", value: ""}, {name: "KNRI11", value: ""}, {name: "MALL11", value: ""}, {name: "MAXR11", value: ""}, {name: "MBRF11", value: ""}, {name: "MFII11", value: ""}, {name: "MGFF11", value: ""}, {name: "MXRF11", value: ""}, {name: "NSLU11", value: ""}, {name: "NVHO11", value: ""}, {name: "ONEF11", value: ""}, {name: "ORPD11", value: ""}, {name: "OUCY11", value: ""}, {name: "OUJP11", value: ""}, {name: "PABY11", value: ""}, {name: "PDDA11B", value: ""}, {name: "PLRI11", value: ""}, {name: "PORD11", value: ""}, {name: "PQDP11", value: ""}, {name: "PRSV11", value: ""}, {name: "RBBV11", value: ""}, {name: "RBCB11", value: ""}, {name: "RBDS11", value: ""}, {name: "RBGS11", value: ""}, {name: "RBRD11", value: ""}, {name: "RBRF11", value: ""}, {name: "RBRR11", value: ""}, {name: "RBVO11", value: ""}, {name: "RDES11", value: ""}, {name: "RDPD11", value: ""}, {name: "REIT11", value: ""}, {name: "RNGO11", value: ""}, {name: "SAAG11", value: ""}, {name: "SCPF11", value: ""}, {name: "SDIL11", value: ""}, {name: "SHPH11", value: ""}, {name: "SPTW11", value: ""}, {name: "TBOF11", value: ""}, {name: "TFOF11", value: ""}, {name: "TGAR11", value: ""}, {name: "THRA11", value: ""}, {name: "TRNT11", value: ""}, {name: "TRXL11", value: ""}, {name: "UBSR11", value: ""}, {name: "VISC11", value: ""}, {name: "VLOL11", value: ""}, {name: "VRTA11", value: ""}, {name: "WPLZ11", value: ""}, {name: "XPCM11", value: ""}, {name: "XPIN11", value: ""}, {name: "XPLG11", value: ""}, {name: "XPML11", value: ""}, {name: "XTED11", value: ""}, {name: "BBSD11", value: ""}, {name: "BOVA11", value: ""}, {name: "BOVV11", value: ""}, {name: "BRAX11", value: ""}, {name: "DIVO11", value: ""}, {name: "ECOO11", value: ""}, {name: "FIND11", value: ""}, {name: "FNAM11", value: ""}, {name: "FNOR11", value: ""}, {name: "FPOR11", value: ""}, {name: "FSRF11", value: ""}, {name: "FSTU11", value: ""}, {name: "ISUS11", value: ""}, {name: "IVVB11", value: ""}, {name: "MATB11", value: ""}, {name: "MMXM11", value: ""}, {name: "PIBB11", value: ""}, {name: "SMAL11", value: ""}, {name: "SPXI11", value: ""}, {name: "XBOV11", value: ""}, {name: "XPOM11", value: ""}, {name: "BEEF11", value: ""}, {name: "FJTA11", value: ""}, {name: "FJTA13", value: ""}, {name: "FJTA15", value: ""}, {name: "FJTA17", value: ""}, {name: "LOGN12", value: ""}, {name: "MYPK11", value: ""}, {name: "MYPK12"}}
