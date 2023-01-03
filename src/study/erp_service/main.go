package main

import (
	"flag"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/schollz/progressbar/v3"
	"github.com/tebeka/selenium"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

var bar *progressbar.ProgressBar

const (
	chromeDriverPath  = "C:\\driver\\96.0.4664.110\\chromedriver"
	port              = 8080
	logisticsBillCode = "发货单号"
)

func main() {
	var username string
	var password string
	flag.StringVar(&username, "u", "", "ERP 账号，必填")
	flag.StringVar(&password, "p", "", "ERP 账号密码，必填")
	flag.Parse()
	username = "****"
	password = "****"
	f, err := excelize.OpenFile("有差异的发货单.xlsx")
	if err != nil {
		log.Fatal(err)
	}
	//读取某个表单的所有数据
	var rows = f.GetRows("Sheet1")
	hasDeal, _ := excelize.OpenFile("./hasDeal.xlsx")
	var rowsHasDeal = hasDeal.GetRows("hasDeal")
	defer hasDeal.SaveAs("./hasDeal.xlsx")
	logisticsBillCodeIndex := -1
	bar = progressbar.Default(int64(len(rows)), "任务进度:")

	// Start a WebDriver server instance
	opts := []selenium.ServiceOption{
		selenium.Output(os.Stderr), // Output debug information to STDERR.
	}
	//selenium.SetDebug(true)
	service, err := selenium.NewChromeDriverService(chromeDriverPath, port, opts...)
	if err != nil {
		panic(err) // panic is used only as an example and is not otherwise recommended.
	}
	defer service.Stop()
	// Connect to the WebDriver instance running locally.
	caps := selenium.Capabilities{"browserName": "chrome"}
	wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", port))
	if err != nil {
		panic(err)
	}
	//defer wd.Quit()
	if err = wd.Get("https://ytd.app.gerpgo.com/auth/login"); err != nil {
		panic(err)
	}

	usernameInput, err := wd.FindElement(selenium.ByCSSSelector, "body > div:nth-child(1) > app-root > div > div > div > app-auth > div > app-login > div > div.center.clearFix > form > nz-form-item:nth-child(2) > nz-form-control > div > span > nz-input-group > input")
	if err != nil {
		//
		//(err)
	}
	usernameInput.SendKeys(username)
	passwordInput, _ := wd.FindElement(selenium.ByCSSSelector, "body > div:nth-child(1) > app-root > div > div > div > app-auth > div > app-login > div > div.center.clearFix > form > nz-form-item:nth-child(3) > nz-form-control > div > span > nz-input-group > input")
	passwordInput.SendKeys(password)
	loginButton, _ := wd.FindElement(selenium.ByCSSSelector, "body > div:nth-child(1) > app-root > div > div > div > app-auth > div > app-login > div > div.center.clearFix > form > nz-form-item:nth-child(5) > nz-form-control > div > span > button")
	loginButton.Click()
	time.Sleep(5 * time.Second)
	var hasDealString = " "
	for i, row := range rowsHasDeal {
		print(i)
		print(row[0])
		print(hasDealString)
		hasDealString = hasDealString + row[0]
	}
	// 数据解析
	for i, row := range rows {
		if logisticsBillCodeIndex < 0 && i == 0 {
			//hasDeal.SetCellValue("hasDeal", "A"+strconv.FormatInt(int64(i), 10), "已处理单号")
			for i, value := range row {
				switch value {
				case logisticsBillCode:
					logisticsBillCodeIndex = i
				}
			}
		}
		print("-------------当前进行到了：%v", i+1)
		code := row[logisticsBillCodeIndex]
		if i != 0 && !strings.Contains(hasDealString, code) {
			if err := wd.Get("https://ytd.app.gerpgo.com/amz-web/logistics/logisticsBill"); err != nil {
				panic(err)
			}
			time.Sleep(3 * time.Second)
			logisticsBillCodeInput, err := wd.FindElement(selenium.ByCSSSelector, "body > div:nth-child(2) > app-root > div > div > div.wrapper.ng-star-inserted > div > app-module-outlet > app-logistics > app-logistics-bill > div > div.list-page-wrapper > div.filters.padding-right > div > div:nth-child(1) > app-batch-dynamic-search > div > app-batch-input-search > nz-input-group > nz-input-group > input")
			if err != nil {
				break
				//panic(err)
			}
			logisticsBillCodeInput.Clear()
			logisticsBillCodeInput.SendKeys(code)
			findByCodeElement, err := wd.FindElement(selenium.ByCSSSelector, "body > div:nth-child(2) > app-root > div > div > div.wrapper.ng-star-inserted > div > app-module-outlet > app-logistics > app-logistics-bill > div > div.list-page-wrapper > div.filters.padding-right > div > div:nth-child(1) > app-batch-dynamic-search > div > app-batch-input-search > nz-input-group > nz-input-group > span > i > svg")
			if err != nil {
				panic(err)
			}
			findByCodeElement.Click()
			time.Sleep(2 * time.Second)
			totalElement, err := wd.FindElement(selenium.ByCSSSelector, "body > div:nth-child(2) > app-root > div > div > div.wrapper.ng-star-inserted > div > app-module-outlet > app-logistics > app-logistics-bill > div > div.list-page-wrapper > div.tab.shadow > div > div.bottom-pagination.fr > nz-pagination > ul > li.ant-pagination-total-text.ng-star-inserted > app-table-pagination-info")
			text, _ := totalElement.Text()
			if strings.EqualFold(strings.TrimSpace(text), "共 0 条记录") {
				continue
			}
			actionElement, _ := wd.FindElement(selenium.ByCSSSelector, "body > div:nth-child(2) > app-root > div > div > div.wrapper.ng-star-inserted > div > app-module-outlet > app-logistics > app-logistics-bill > div > div.list-page-wrapper > div.tab.shadow > nz-table > nz-spin > div > div > div > div > div.ant-table-body.ng-star-inserted > table > tbody > tr > td.ant-table-td-right-sticky > div > nz-dropdown-button > div > button.ant-dropdown-trigger.ant-btn.ant-btn-default.ant-btn-icon-only")
			actionElement.Click()
			time.Sleep(1 * time.Second)
			//modalOpenElement, _ := wd.FindElement(selenium.ByCSSSelector, "#cdk-overlay-7 > div > ul > li:nth-child(2)")
			//ulElement, _ := wd.FindElement(selenium.ByID, "cdk-overlay-9")
			liElements, _ := wd.FindElements(selenium.ByTagName, "li")
			for _, liElement := range liElements {
				liText, _ := liElement.Text()
				if strings.EqualFold(liText, "处理收发异常") {
					liElement.Click()
					time.Sleep(1 * time.Second)
					finishDateInputElement, _ := wd.FindElement(selenium.ByCSSSelector, "body > div:nth-child(2) > app-root > div > div > div.wrapper.ng-star-inserted > div > app-module-outlet > app-logistics > app-logistics-bill > div > div.except-header-container-mount > nz-modal > div > div.ant-modal-wrap.vertical-center-modal.full-plus > div > div > div.ant-modal-body.ng-star-inserted > app-logistics-bill-finished > nz-spin > div > div > form > div.ant-row-flex.ant-row-flex-top.ant-row-flex-start > div.text-ellipsis-2.ant-col-6 > nz-form-item > nz-form-control > div > span > nz-date-picker > nz-picker > span > input")
					finishDateInputElement.Click()
					time.Sleep(2 * time.Second)
					finishDateButtonElement, _ := wd.FindElement(selenium.ByCSSSelector, "#cdk-overlay-7 > div > date-range-popup > div > div > div > div > inner-popup > div > date-table > table > tbody > tr.ant-calendar-current-week.ant-calendar-active-week.ng-star-inserted > td.ant-calendar-cell.ant-calendar-today.ant-calendar-selected-day.ng-star-inserted > div")
					time.Sleep(1 * time.Second)
					finishDateButtonElement.Click()
					//9-12
					tableElement, _ := wd.FindElement(selenium.ByCSSSelector, "body > div:nth-child(2) > app-root > div > div > div.wrapper.ng-star-inserted > div > app-module-outlet > app-logistics > app-logistics-bill > div > div.except-header-container-mount > nz-modal > div > div.ant-modal-wrap.vertical-center-modal.full-plus > div > div > div.ant-modal-body.ng-star-inserted > app-logistics-bill-finished > nz-spin > div > div > form > div.ant-row > div > div > nz-table")
					trElements, _ := tableElement.FindElements(selenium.ByTagName, "tr")
					for _, trElement := range trElements {
						//body > div:nth-child(2) > app-root > div > div > div.wrapper.ng-star-inserted > div > app-module-outlet > app-logistics > app-logistics-bill > div > div.except-header-container-mount > nz-modal > div > div.ant-modal-wrap.vertical-center-modal.full-plus > div > div > div.ant-modal-body.ng-star-inserted > app-logistics-bill-finished > nz-spin > div > div > form > div.ant-row > div > div > nz-table > nz-spin > div > div > div > div > div.ant-table-body.ng-star-inserted > table > tbody > tr.gj-tr.gj-tr-has-picture.tabTr.ng-untouched.ng-pristine.ng-valid.ant-table-row.ng-star-inserted
						tdElements, _ := trElement.FindElements(selenium.ByTagName, "td")
						if len(tdElements) > 0 {
							differenceString, err := tdElements[9].Text()
							if err != nil {
								//panic(err)
							}
							difference, err := strconv.Atoi(differenceString)
							if err != nil {
								continue
							}
							differenceabs := math.Abs(float64(difference))
							y := int64(differenceabs)
							tiaozhengElement, err := tdElements[12].FindElement(selenium.ByTagName, "input")
							if err != nil {
								break
								//panic(err)
							}
							tiaozhengElement.Clear()
							tiaozhengElement.SendKeys(strconv.FormatInt(y, 10))
						}
					}
					time.Sleep(1 * time.Second)
					submitElement, err := wd.FindElement(selenium.ByCSSSelector, "body > div:nth-child(2) > app-root > div > div > div.wrapper.ng-star-inserted > div > app-module-outlet > app-logistics > app-logistics-bill > div > div.except-header-container-mount > nz-modal > div > div.ant-modal-wrap.vertical-center-modal.full-plus > div > div > div.ant-modal-footer.ng-star-inserted > button.ant-btn.ng-star-inserted.ant-btn-default.ant-btn-primary")
					if err != nil {
						break
						//panic(err)
					}
					submitElement.Click()
					break
				}
			}
			bar.Add(1)
			hasDeal.SetCellValue("hasDeal", "A"+strconv.FormatInt(int64(i+len(rowsHasDeal)), 10), code)
		}
	}

}
