package mail

import (
	"gopkg.in/gomail.v2"
	"strconv"
)

// send mail
func SendMail(mailTo []string, subject string, body string) error {
	//定义邮箱服务器连接信息，如果是阿里邮箱 pass填密码，qq邮箱填授权码
	mailConn := map[string]string{
		"user": "imagesharing@163.com",
		"pass": "DSKJFSZUYCREFRFU",
		"host": "smtp.163.com",
		"port": "25",
	}

	port, _ := strconv.Atoi(mailConn["port"]) //转换端口类型为int

	m := gomail.NewMessage()
	m.SetHeader("From", mailConn["user"]) //这种方式可以添加别名，即“XD Game”， 也可以直接用<code>m.SetHeader("From",mailConn["user"])</code> 读者可以自行实验下效果
	m.SetHeader("To", mailTo...)          //发送给多个用户
	m.SetHeader("Subject", subject)       //设置邮件主题
	m.SetBody("text/html", body)          //设置邮件正文

	d := gomail.NewDialer(mailConn["host"], port, mailConn["user"], mailConn["pass"])

	err := d.DialAndSend(m)
	return err

}
func main() {
	//定义收件人
	mailTo := []string{
		"2296176046@qq.com",
	}
	//邮件主题为"Hello"
	subject := "Hello"
	// 邮件正文
	body := "Good"
	SendMail(mailTo, subject, body)
}
