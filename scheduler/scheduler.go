package scheduler

import (
	"github.com/cihub/seelog"
	"os"
)

// NewLogger will load the scheduler's configuration file.
// If file name is not supplied, it will use default configuration.
func SetLogger(fileName string) {
	if _, err := os.Stat(fileName); err == nil {
		logger, err := seelog.LoggerFromConfigAsFile(fileName)
		if err != nil {
			panic(err)
		}

		seelog.ReplaceLogger(logger)
	} else {
		configString := `<scheduler>
                        <outputs formatid="main">
                            <filter levels="info,error,critical">
                                <rollingfile type="date" filename="log/AppLog.log" namemode="prefix" datepattern="02.01.2006"/>
                            </filter>
                            <console/>
                        </outputs>
                        <formats>
                            <format id="main" format="%Date %Time [%LEVEL] %Msg%n"/>
                        </formats>
                        </scheduler>`
		logger, err := seelog.LoggerFromConfigAsString(configString)
		if err != nil {
			panic(err)
		}

		seelog.ReplaceLogger(logger)
	}

	seelog.Info("log initialize finish.")
}