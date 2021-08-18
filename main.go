/*
扑克牌52张，花色黑桃spades，红心hearts，方块diamonds，草花clubs各13张，2-10，J，Q，K，A

Face：即2-10，J，Q，K，A，其中10用T来表示。

Color：即S(spades)、H(hearts)、D(diamonds)、C(clubs)

用 Face字母+小写Color字母表示一张牌，比如As表示黑桃A，其中A为牌面，s为spades，即黑桃，Ah即红心A，以此类推。
现分别给定任意两手牌各有5张，例如：AsAhQsQhQc，vs KsKhKdKc2c，请按德州扑克的大小规则来判断双方大小。

代码要求有通用性，可以任意输入一手牌或几手牌来进行比较。


5张
*/

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Poker struct {
	Alice string `json:"alice"`
	Bob string	`json:"bob"`
	Result int	`json:"result"`
}

type Match struct {
	Matches []Poker `json:"matches"`
}

// 定义一个牌型
// 1.皇家同花顺
// 2.同花顺
// 3.四条
// 4.满堂彩（葫芦，三带二）
// 5.同花
// 6.顺子
// 7.三条
// 8.两对
// 9.一对
// 10.单张大牌
type CardType int

// 把数据读取出来 分别放在切片中
func ReadFile(filename string) (alices,bobs []string, results []int){
	buf ,err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	var matches Match
	err = json.Unmarshal(buf,&matches)
	if err != nil {
		panic(err)
	}


	alices = make([]string,10240)
	bobs = make([]string,10240)
	results = make([]int,10240)

	for k, v := range matches.Matches {
		//fmt.Printf("k=%#v,v=%#v\n",k,v)
		alices[k] = v.Alice
		bobs[k] = v.Bob
		results[k] = v.Result
	}
	return
}

// 根据手牌判断 牌型
func JudgmentGroup(cards string)(cardType CardType,cardSizes, cardColors []string ) {
	// 遍历分开手牌的大小和颜色
	cardSizes = make([]string,5)
	cardColors = make([]string,5)
	for i,v := range cards {
		// fmt.Println("v=",string(v))
		if i%2 !=0 {
			cardColors[i/2] = string(v)
		}else {
			cardSizes[i/2] = string(v)
		}
	}

	// 先判断是否是同花
	var s, h, d,c int
	for _, v := range cardColors {

		switch v {
		case "h": // 红心
			h++
		case "s": // 黑桃
			s++
		case "d": // 方块
			d++
		case "c": // 草花
			c++
		default:
			fmt.Println("无法解析的花色")
			os.Exit(-1)
		}
	}
	// 同花
	if s == 5 || h == 5 || d == 5 || c == 5 {
		// 同花后判断是不是顺子

		// 判断是不是顺子
		isShun := IsShunZi(cardSizes)

		if isShun == true {
			// 是 就同花顺
			cardType = 1
		}else {
			// 否 就同花
			cardType = 5
		}
		return

	}

	// 不是同花
	// 先判断是不是顺子
	isShun := IsShunZi(cardSizes)
	if isShun == true {
		// 顺子出函数
		cardType = 6
		return
	}
	//  判断重复张数
	// var sum []int
	sums := make([]int,len(cardSizes)/2 +1)
	for i := 0; i<len(cardSizes)/2 +1 ; i++ {
		for _,v := range cardSizes {
			//fmt.Printf("v111=%#v--->",v)
			if v == cardSizes[i] {
				sums[i] ++
			}
		}
		//fmt.Println("sum=",sums[i])
	}
	k1 := 0 // 记录三个
	k2 := 0 // 记录两个
	for _,v := range sums {
		if v >=4 {
			// 为四条出函数
			cardType = 3
			return
		}
		if v == 3 {
			k1++
		}
		if v == 2 {
			k2++
		}
	}
	//fmt.Printf("k1=%#v,k2=%#v\n",k1,k2)
	// 三带二 出函数
	if k1 >=1 && k2 >=1 {
		cardType = 4
		return
	}
	// 三条
	if k1 >= 1 {
		cardType = 7
		return
	}
	// 两对
	if k2 >= 2 {
		cardType = 8
		return
	}
	// 一对
	for i:=len(cardSizes)/2 +1 ; i< len(cardSizes);i++ {
		for _,v := range cardSizes {
			//fmt.Printf("v111=%#v--->",v)
			if v == cardSizes[i] {
				sums[i/2] ++
			}
		}
		//fmt.Println("sum=",sums[i/2])
	}
	for _,v := range sums {
		if v >2 {
			cardType = 9
			return
		}
	}
	// 单张大牌
	cardType=10
	return
}

// 判断是不是顺子
func IsShunZi(cardSize []string)(shunZi bool){
	shunZi = false
	saves := make([]string,13)
	for _, v := range cardSize {
		switch v {
		case "2":
			saves[0] = v
		case "3":
			saves[1]= v
		case "4":
			saves[2] = v
		case "5":
			saves[3] = v
		case "6":
			saves[4] =v
		case "7":
			saves[5]=v
		case "8":
			saves[6]=v
		case "9":
			saves[7] =v
		case "T":
			saves[8] = v
		case "J":
			saves[9]=v
		case "Q":
			saves[10]=v
		case "K":
			saves[11]=v
		case "A":
			saves[12] = v
		default:
			fmt.Println("无法解析的扑克牌")
		}
	}
	// 判断数组是否连续
	sum := 0
	for _ ,v := range saves {
		if v != "" {
			sum ++
		}else {
			sum = 0
		}
		if sum >=5 {
			// break
			shunZi = true
			return
		}
	}
	return
}

// 下面是同类型的比较 ---------

// 面值转换
func tranNumFace(num uint8) (f32 float32) {
	switch num {
	case 48 :
		f32 = 0
	case 50 :
		f32 = 2
	case 51:
		f32 = 3
	case 52:
		f32 = 4
	case 53:
		f32 = 5
	case 54:
		f32 = 6
	case 55:
		f32 = 7
	case 56:
		f32 = 8
	case 57:
		f32 = 9
	case 84:
		f32 = 10
	case 74:
		f32 = 11
	case 81:
		f32 = 12
	case 75:
		f32 = 13
	case 65:
		f32 = 14
	default:
		fmt.Println("无法解析的花色---退出程序")
		os.Exit(-1)
	}
	return
}
// 花色转换
func tranNumColor(num uint8) (f32 float32) {
	switch num {
	case 100:
		f32 = 0.0
	case 99:
		f32 = 0.1
	case 104:
		f32 = 0.2
	case 115:
		f32 = 0.3
	}
	return
}

// 牌string转成float32切片,并且返回一个有顺序的
func TranNums(card string ) (f32s []float32) {
	//转换
	f32s = make([]float32,len(card)/2)
	for i:=0; i< len(card); i+=2 {
		f32s[i/2] = tranNumFace(card[i])/*+tranNumColor(card[i+1])*/
	}
	//排序
	//fmt.Println("f32s---->",f32s)
	f32s = QuickSortFloat32(f32s)
	//fmt.Println("after--f32s->",f32s)
	return
}

// 快排  -->倒序
func QuickSortFloat32(f32s []float32) []float32 {
	if len(f32s) <= 1 {
		return f32s
	}
	splitdata := f32s[0]          //第一个数据
	low := make([]float32, 0, 0)     //比我小的数据
	hight := make([]float32, 0, 0)   //比我大的数据
	mid := make([]float32, 0, 0)     //与我一样大的数据
	mid = append(mid, splitdata) //加入一个
	for i := 1; i < len(f32s); i++ {
		if f32s[i] > splitdata {
			low = append(low, f32s[i])
		} else if f32s[i] < splitdata {
			hight = append(hight, f32s[i])
		} else {
			mid = append(mid, f32s[i])
		}
	}
	low, hight = QuickSortFloat32(low), QuickSortFloat32(hight)
	myarr := append(append(low, mid...), hight...)
	return myarr
}

// 单张大牌的比较  -----1前大 2后大  0一样大 -1出错
func SingleCardCompareSize(cards ...string)(result int) {
	//fmt.Println("=========>len(cards)=",len(cards))
	mapf32s := make(map[int][]float32,len(cards))

	// 分别把传进来的很多副手牌各自排序好-->[]float32
	for i, v := range cards {
		//fmt.Printf("i=%#v,v=%#v\n",i,v)
		mapf32s[i] = make([]float32,len(v))
		mapf32s[i] = TranNums(v)
	}
	if len(cards) == 2 {
		for i :=0; i<len(mapf32s[0]); i++ {
			if mapf32s[0][i] >mapf32s[1][i] {
				return 1
			}else if mapf32s[0][i] < mapf32s[1][i]{
				return 2
			}
		}
		return 0
	}
	// 比较两幅以上....
	return -1
}

// 一对比大小
/*则两张牌中点数大的赢，如果对牌都一样，
则比较另外三张牌中大的赢，如果另外三张牌中较大的也一样则比较第二大的和第三大的，如果所有的牌都一样，则平分彩池。*/
func aPairCom(cardLen int,card ... string)(result int) {
	// 先找出对子
	pairs := make([]uint8,cardLen)
	for index, _ := range card {
		fmt.Printf("--->begin----card[%#v]=%#v\n",index,card[index])
		b := []byte(card[index])
		for i:=0; i< cardLen/2+2; i+=2 {
			for j := i+2; j<cardLen;j+=2 {
				if b[i] == b[j] {
					// 获取对子
					pairs[index] = card[index][i]
					fmt.Println("pairs[index]=",pairs[index])
					b[i] = 48  // 0
					b[j] = 48
				}
			}
		}
		card[index] = string(b)
		fmt.Printf("--->after----card[%#v]=%#v\n",index,card[index])
	}
	// 传进来的参数两个
	if len(card) == 2 {
		if pairs[0] > pairs[1] {
			result = 1
			return
		}else if pairs[0] < pairs[1] {
			result = 2
			return
		}else {
			// 对子一样，比较剩余牌的值大小
			// 把对子移除
			result = SingleCardCompareSize(card[0],card[1])
			return
		}
	}

	return
}

// 两对
/*两对
两对点数相同但两两不同的扑克和随意的一张牌组成。
平手牌：如果不止一人抓大此牌相，牌点比较大的人赢，如果比较大的牌点相同，那么较小牌点中的较大者赢，
如果两对牌点相同，那么第五张牌点较大者赢（起脚牌）。如果起脚牌也相同，则平分彩池。*/
func twoPairCom(cardLen int, card ... string)(result int) {

}

func PokerMan() {
	file := "/home/weilijie/chromeDown/match_result.json"
	alices := make([]string,1024)
	bobs := make([]string,1024)
	results := make([]int,1024)
	alices,bobs,results = ReadFile(file)
	//return

	// for i:=0; i < len(alices); i ++ {
	// 	fmt.Printf("alices[%#v]=%#v\n",i,alices[i])
	// 	fmt.Printf("bobs[%#v]=%#v\n",i,bobs[i])
	// 	fmt.Printf("results[%#v]=%#v\n",i,results[i])
	// }

	for i:=0; i < len(alices); i ++ {
		result := -1
		val1,_,_ := JudgmentGroup(alices[i])
		val2,_,_ := JudgmentGroup(bobs[i])
		if val1 < val2 {
			result = 1
		}else if val1 > val2 {
			result = 2
		}else {
			// 牌型处理相同的情况
			// ...

			switch val1 {
			case 10:
				// 同类型下的单张大牌比较
				result = SingleCardCompareSize(alices[i],bobs[i])
			case 9:
				// 同类型的一对
				result = aPairCom(10,alices[i],bobs[i])

			}

			// 最后比较结果
		}

		if result != results[i] {
			fmt.Printf("判断错误--->alice:%#v,bob:%#v<-----\n",alices[i],bobs[i])
		} else {
			fmt.Println("判断正确222222")
		}
	}
}


func main() {
	// file := "/home/weilijie/chromeDown/match_result.json"
	// ReadFile(file)

	// JudgmentGroup("AsAhQsQhQc")
	// str1 := "23456789"
	// fmt.Println("[]byte(str1)",[]byte(str1))
	// str2 := "T"
	// fmt.Println("[]byte(str2)",[]byte(str2))
	// str3 := "JQKA"
	// fmt.Println("[]byte(str3)",[]byte(str3))
	// str := []string {"2"}
	// isres := IsShunZi(str)
	// fmt.Println("isres=",isres)
	//
	// type1 ,_,_ := JudgmentGroup("KsKhKdKc2c")
	// fmt.Println("type1=",type1)

	// result := SingleCardCompareSize(5,"AhKdKc2cKs","2cKsKhKdKc")
	// fmt.Println("---->result=",result)
	// str4 := "shcd"
	// fmt.Println("[]byte(str4)=",[]byte(str4))
	//
	// var int32 int32
	// int32 = 105*1000
	// fmt.Println("int_32=",int32)
	// fs :=[]float32 {88.88,11.0,99.2,1.0,5.4,88.2,77,65,20,3}
	// res := make([]float32,10)
	// res =QuickSortFloat32(fs)
	// fmt.Printf("fs=%#v\nres=%#v\n",fs,res)


	//PokerMan()

	// str := "0"
	// fmt.Println("str=",[]byte(str))
	result := aPairCom(10,"2d4s4h3s5h","4h2h5s3c4d")
	fmt.Println("result=",result)
}
