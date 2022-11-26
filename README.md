# refactoring
refactoring - это тестовое упражнение для компании "TrueConf".
# Предложения по оптимизации
## string to bytes
Главным преимуществом и главным недостатком Golang является GCC. Он позволяет не задумываться о ручном управлении памяти, однако не в меньшей мере влияет на производительность приложений. Поэтому в Golang разработке принято минимизировать использование кучи. Конкретно в этом примере, я бы мог предложить оnкзатся от записи типа
```
[]byte(string)
```
в пользу 2 методов Str2bts и Bts2str определённых в отдельном пакете в директории /pkg. Это снизит использование кучи и повысит производительность.
```golang

func Bts2str(b []byte) (s string) {
	return *(*string)(unsafe.Pointer(&b))
}

func Str2bts(s string) (b []byte) {
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))

	bh.Data = sh.Data
	bh.Cap = sh.Len
	bh.Len = sh.Len

	return b
}
```
## Decode/Encode
Следующим пунктом я хотел бы отметить использование стандартных библиотек для кодирования и декодирования JSON. Существует множество библиотек выполняющих эту задачу
эффективнее как по используемой памяти, так и по скорости. Вы можете сравнить все и выбрать то, что будет удобно:
1. [json-iterator/go](https://github.com/json-iterator/go)
2. [easyjson](https://github.com/mailru/easyjson)
3. [gojay](https://github.com/francoispqt/gojay)
4. [segmentio/encoding/json](https://github.com/segmentio/encoding)
5. [jettison](https://github.com/wI2L/jettison)
6. [simdjson-go](https://github.com/minio/simdjson-go)

Также отмечу, что JSON не едиственных формат обмена данными. Альтернативами ему могут служить GOB и Protobuf.
## Zero allocation and Non blocking
Стандартная библиотека net/http вынуждена на каждый запрос запускать горутину тем самым используя большое количество памяти, а также предоставляет блокирующий
API-интерфейс ввод-вывода. Существуют библиотеки исправлющие эти недостатки. [Hertz](https://github.com/cloudwego/hertz) (если нужен HTTP) и
[Kitex](https://github.com/cloudwego/kitex) (если нужен RPC) - это 2 наиболее приемлемых фреймворка для
использования в продакшене, созданные компанией ByteDance.

