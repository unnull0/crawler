package main

import (
	"fmt"
	"time"

	"github.com/unnull0/crawler/grabber"
	"github.com/unnull0/crawler/grabber/workerengine"
	"github.com/unnull0/crawler/log"
	"github.com/unnull0/crawler/tasklib/doubantenement"
	"go.uber.org/zap/zapcore"
)

func main() {
	core, c := log.NewFileCore("./text.log", zapcore.InfoLevel)
	defer c.Close()
	logger := log.NewLogger(core)
	logger.Info("log init end")

	cookie := `douban-fav-remind=1; viewed="1007305"; bid=-LBsoZdQj8k; gr_user_id=9381f0c8-c9a0-4ef7-9e10-1b4812fb785e; __gads=ID=a486bca57b09d9fc-224d131d48d900aa:T=1673615555:RT=1673615555:S=ALNI_MZvaEzveWJzMynGcfoXl7iHp8T1WA; __yadk_uid=fLQE8dmo4P1yoZDK9bOwSzQusWeSs75z; push_noty_num=0; push_doumail_num=0; __utmv=30149280.26736; ap_v=0,6.0; __gpi=UID=00000ba35654ce16:T=1673615555:RT=1678430610:S=ALNI_Mb-V6rgnowh8JznMa81jye-o5mh8g; apiKey=; dbcl2="267368696:rLWKgwWt9MY"; ck=k7w3; _pk_ref.100001.8cb4=["","",1678433761,"https://accounts.douban.com/"]; _pk_id.100001.8cb4=55cd98221bb9ae38.1612278505.9.1678433761.1678430715.; _pk_ses.100001.8cb4=*; __utma=30149280.1521688018.1592847100.1678430612.1678433762.13; __utmc=30149280; __utmz=30149280.1678433762.13.4.utmcsr=accounts.douban.com|utmccn=(referral)|utmcmd=referral|utmcct=/; __utmt=1; __utmb=30149280.7.5.1678433762`

	fetcher := grabber.NewFetcher(grabber.BrowserFetchType)

	var workerlist []*grabber.Request
	for i := 0; i <= 25; i += 25 {
		str := fmt.Sprintf("https://www.douban.com/group/szsh/discussion?start=%d", i)
		workerlist = append(workerlist, &grabber.Request{
			URL:       str,
			Cookie:    cookie,
			Timeout:   3000 * time.Millisecond,
			WaitTime:  1 * time.Second,
			ParseFunc: doubantenement.ParseURL,
		})
	}

	s := workerengine.NewSchedule(
		workerengine.WithFetcher(fetcher),
		workerengine.WithLogger(logger),
		workerengine.WithSeeds(workerlist),
		workerengine.WithWorkCount(5),
	)
	s.Run()
}
