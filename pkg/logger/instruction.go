package logger

import "fmt"

func main() {
	/*
		DebugLevel：	调试级别，通常在生产环境中禁用，日志量通常较大。
		InfoLevel：		信息级别，是默认的日志优先级。
		WarnLevel：		警告级别，比信息级别更重要，但不需要人工逐个审查。
		ErrorLevel：	错误级别，高优先级。如果应用程序正常运行，不应生成任何错误级别的日志。
		DPanicLevel：	重要错误级别，开发环境中的日志记录器在写入消息后会导致 panic。
		PanicLevel：	紧急错误级别，记录一条消息后立即引发 panic。
		FatalLevel：	严重错误级别，记录一条消息后调用 os.Exit(1) 终止程序。
	*/

	/*
		1.对于这段代码：

		// 创建一个用于终端输出的核心，并将其添加到 cores 数组中。
		if level == "debug" { // 如果是debug级别好需要输出到终端
			// 创建终端输出编码器
			consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
			cores = append(cores, zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), zapcore.DebugLevel))
		}

		解释：
			如果 level 参数等于 "debug"，则创建一个用于终端输出的核心，并将其添加到 cores 数组中。
			该核心使用 zapcore.NewConsoleEncoder 创建终端输出编码器，zapcore.Lock(os.Stdout) 锁定标准输出流，以及 zapcore.DebugLevel 日志级别。

			在给定的代码段中，zapcore.Lock(os.Stdout) 表示锁定标准输出流。让我们来详细解释这个概念。
				在 Go 语言中，os.Stdout 是一个标准的输出流，通常指向终端控制台或命令行界面的标准输出。
			通过将 os.Stdout 作为参数传递给 zapcore.Lock 函数，可以锁定标准输出流，以确保在并发环境中不会发生交叉输出或竞争条件。

			锁定标准输出流的主要目的是避免多个日志线程同时尝试写入控制台输出时产生的混乱或交错的输出。
			通过对标准输出流进行加锁，每个线程将依次写入日志消息，从而保证输出的顺序性和一致性。

			因此，通过在 cores 数组中添加具有锁定标准输出流的核心，可以将日志消息输出到终端控制台，并确保在多线程或并发情况下输出的正确性。
			这样可以方便地在终端上查看日志输出，特别是在调试过程中或需要即时查看日志的情况下。

	*/
	fmt.Println("--------------------------------------------------------------------------------------")
	/*
		2.对于这段代码
			return &Log{zap.New(core, zap.AddCaller())}

		这里的 core 是之前创建的日志核心，它是通过多个核心（cores 数组）组合而成的。core 包含了多种不同的日志输出方式和级别设置。

		zap.New 函数是日志库提供的用于创建日志记录器的方法。它接受一个日志核心和一些可选的配置选项，并返回一个新的日志记录器。

		在这段代码中，使用 core 作为日志核心参数传递给 zap.New 函数，然后使用 zap.AddCaller() 选项将调用者的信息添加到日志中。

		添加调用者的信息到日志中意味着在日志输出中包含记录该日志消息的代码位置（调用者的文件名、函数名和行号）。
		这对于定位日志消息的来源非常有用，特别是在调试和排查问题时。通过将 zap.AddCaller() 选项传递给 zap.New 函数，日志记录器会自动捕获和包含每条日志消息的调用者信息。

		因此，通过这段代码，创建了一个新的日志记录器，并将其作为指针返回给调用方。
		该日志记录器将使用之前定义的日志核心，并且在日志输出中包含调用者的信息。
	*/
	fmt.Println("--------------------------------------------------------------------------------------")

	/*
		对于这段代码 getEncoder()
		func getEncoder() zapcore.Encoder {
			encoderConfig := zapcore.EncoderConfig{
				TimeKey:        "ts",
				LevelKey:       "level",
				NameKey:        "logger",
				CallerKey:      "caller",
				FunctionKey:    zapcore.OmitKey,
				MessageKey:     "msg",
				StacktraceKey:  "stacktrace",
				LineEnding:     zapcore.DefaultLineEnding,   // 结尾字符
				EncodeLevel:    zapcore.CapitalLevelEncoder, // 将Level序列化为全大写字符串。例如, InfoLevel被序列化为INFO
				EncodeTime:     zapcore.ISO8601TimeEncoder,  // 格式化时间戳
				EncodeDuration: zapcore.SecondsDurationEncoder,
				EncodeCaller:   zapcore.ShortCallerEncoder,
			}
			return zapcore.NewConsoleEncoder(encoderConfig) // 普通的Log Encoder
		}

		TimeKey：		时间戳字段的键名。
		LevelKey：		日志级别字段的键名。
		NameKey：		日志记录器名称字段的键名。
		CallerKey：		调用者信息字段的键名。
		FunctionKey：	函数名字段的键名。
		MessageKey： 	日志消息字段的键名。
		StacktraceKey：	堆栈跟踪信息字段的键名。
		LineEnding：	行尾字符，使用 zapcore.DefaultLineEnding 表示默认的行尾字符。
		EncodeLevel：	日志级别的编码方式，使用 zapcore.CapitalLevelEncoder 表示将级别序列化为全大写字符串。
		EncodeTime：	时间戳的编码方式，使用 zapcore.ISO8601TimeEncoder 表示将时间戳格式化为 ISO8601 格式。
		EncodeDuration：持续时间的编码方式，使用 zapcore.SecondsDurationEncoder 表示将持续时间编码为秒数。
		EncodeCaller：	调用者信息的编码方式，使用 zapcore.ShortCallerEncoder 表示将调用者信息编码为短格式。

	*/
}
