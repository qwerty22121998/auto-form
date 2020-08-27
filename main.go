package main

import (
	"fmt"
	"github.com/tebeka/selenium"
	"log"
	"math/rand"
	os "os"
	"time"
)

const (
	seleniumPath    = "selenium-server.jar"
	geckoDriverPath = "geckodriver"
	formUrl         = "https://docs.google.com/forms/d/e/1FAIpQLSfDCuDusZhSJ1R9UouwlW4JEjL0IlSaNpoPMs__-VdAZOhBqA/viewform"
	port            = 8080
	script          = `let getRandomInt = (min, max) => {
    min = Math.ceil(min);
    max = Math.floor(max);
    return Math.floor(Math.random() * (max - min + 1)) + min;
};

let questionType = [
    {
        root: "freebirdFormviewerComponentsQuestionRadioRoot",
        ele: "appsMaterialWizToggleRadiogroupOffRadio",
    },
    {
        root: "appsMaterialWizToggleRadiogroupGroupContent",
        ele: "appsMaterialWizToggleRadiogroupRadioButtonContainer",
    },
];

let process = new Promise(resolve => {
    questionType.forEach(type => {
        let root = document.getElementsByClassName(type.root);
        Array.from(root).forEach(e => {
            let choices = e.getElementsByClassName(type.ele);
            let index = getRandomInt(0, choices.length - 1);
            choices[index].click();
        });
    });
    resolve(true);
})
    .then(_ => {
        document
            .getElementsByClassName("freebirdFormviewerViewNavigationSubmitButton")[0]
            .click();
    })
    .catch(e => console.error(e))`
)

func handle(err error) {
	if err != nil {
		panic(err)
	}
}

type Distribute []int

type Question struct {
	Root       string
	Ele        string
	Distribute []Distribute
}

var question = Question{
	Root: ".appsMaterialWizToggleRadiogroupGroupContent",
	Ele:  ".appsMaterialWizToggleRadiogroupRadioButtonContainer",
	Distribute: []Distribute{
		[]int{20, 60, 60, 60, 0},  // Bạn là sinh viên khóa nào? *
		[]int{150, 50},            // Bạn có gửi xe trong trường không? *
		[]int{80, 80, 40},         // Bạn hiện đang cư trú tại đâu? *
		[]int{50, 80, 50, 15, 5},  // Mức độ hài lòng của bạn về dịch vụ gửi xe tại trường
		[]int{20, 40, 80, 40, 20}, // Bạn cảm thấy bất tiện nhất về vấn đề nào khi sử dụng hình thức ghi vé xe bằng giấy?
		[]int{180, 20},            // Theo bạn, phương thức gửi xe bằng thẻ có góp phần cải thiện dịch vụ nhà xe hay không? *
		[]int{40, 40, 100, 20},    // Bạn cảm thấy bất tiện nhất về vấn đề nào khi nhà xe hiện tại ở xa so với khu vực dom E, F, G, H và toà Gamma? *
		[]int{5, 10, 25, 40, 120}, // Theo bạn, mức độ cần thiết của việc mở rộng nhà xe hiện tại là như thế nào? *
		[]int{10, 10, 60, 60, 60}, // Phân luồng các phương tiện lưu thông trong nhà xe
		[]int{10, 10, 60, 60, 60}, // Tách cổng ra cổng vào riêng biệt
		[]int{10, 10, 60, 60, 60}, // Mở thêm đường cho người đi bộ ra vào lấy xe
		[]int{10, 10, 60, 60, 60}, // Mở rộng khu vực ra vào
		[]int{0, 0, 30, 100, 70},  // Cung cấp dịch vụ cho mượn mũ bảo hiểm, áo mưa
		[]int{10, 10, 60, 60, 60}, // Tuyên truyền nâng cao ý thức của sinh viên
		[]int{10, 10, 60, 60, 60}, // Lắp thêm camera an ninh
		[]int{10, 10, 60, 60, 60}, // Tăng nặng hình thức xử phạt nếu vi phạm
		[]int{10, 10, 60, 60, 60}, // Tăng cường lực lượng an ninh
		[]int{10, 10, 40, 70, 70}, // Ghi mã số sinh viên kèm theo chữ ký
		[]int{10, 10, 40, 70, 70}, // Để lại chứng minh thư, thẻ sinh viên
		[]int{10, 10, 40, 70, 70}, // Giới hạn thời gian cho mượn
		[]int{10, 10, 40, 70, 70}, // Có hình thức đền bù nếu làm hỏng, mất
		[]int{10, 10, 40, 70, 70}, // Có hình thức xử phạt nếu không tuân thủ theo quy định (mượn quá thời gian,...)
		[]int{20, 180},            // Bạn đã bao giờ được phổ biến về các quy định và thủ tục trong việc gửi xe chưa? *
		[]int{0, 0, 30, 90, 80},   // Tổ chức các buổi tập huấn
		[]int{10, 10, 60, 60, 60}, // Phổ biến qua các kênh thông tin của trường
		[]int{10, 10, 60, 60, 60}, // Gửi mail cho sinh viên lúc mới nhập học
		[]int{10, 10, 60, 60, 60}, // Phát hành sổ tay cá nhân
		[]int{180, 20},            // Bạn có cảm thấy cần thiết lập một nhóm mới để mọi người có thể trao đổi, giúp đỡ nhau các vấn đề nói trên? *
	},
}

func initialize() {
	rand.Seed(time.Now().UnixNano())
	for i, v := range question.Distribute {
		var arr []int
		for num, count := range v {
			for count != 0 {
				count--
				arr = append(arr, num)
			}
		}
		rand.Shuffle(len(arr), func(i, j int) {
			arr[i], arr[j] = arr[j], arr[i]
		})
		question.Distribute[i] = arr
	}

}

func fill(wd selenium.WebDriver) {
	handle(wd.Get(formUrl))

	roots, _ := wd.FindElements(selenium.ByCSSSelector, question.Root)

	for i, v := range roots {
		v = v
		choices, _ := v.FindElements(selenium.ByCSSSelector, question.Ele)
		index := question.Distribute[i][0]
		question.Distribute[i] = question.Distribute[i][1:]
		choices[index].Click()
	}
}

func main() {
	initialize()
	opts := []selenium.ServiceOption{
		selenium.GeckoDriver(geckoDriverPath),
		selenium.Output(os.Stderr),
	}
	selenium.SetDebug(true)

	service, err := selenium.NewSeleniumService(seleniumPath, port, opts...)
	handle(err)

	defer service.Stop()
	caps := selenium.Capabilities{"browserName": "firefox"}

	wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", port))
	handle(err)
	defer wd.Quit()
	max := 200
	now := 0
	for now < max {
		now++
		log.Println(now, "/", max)
		fill(wd)
	}

}
