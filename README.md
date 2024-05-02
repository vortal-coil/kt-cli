# ktCloud CLI Client

This is a command-line interface (CLI) client for the ktCloud service. It is written in Go and provides peer-to-peer encryption and zero-trust security. The client allows you to interact with the ktCloud service, enabling you to download, upload, and use the API for ktCloud.

## Features

- **P2P Encryption**: All data transferred between the client and ktCloud service is encrypted using peer-to-peer encryption, ensuring the security of your data.
- **Zero-Trust Security**: The client implements a zero-trust security model, meaning it does not inherently trust any entities. This reduces the potential attack surface.
- **Download and Upload**: The client allows you to download and upload files to and from the ktCloud service.
- **API Interaction**: The client provides a way to interact with the ktCloud API, allowing you to perform various operations on the ktCloud service.

## Installation

To get the ktCloud CLI client or its libraries, you need to have Go installed on your machine. Once Go is installed, you can download and install the ktCloud CLI client using the `go get` command:

```bash
go get github.com/kt-soft-dev/kt-cli
```

For ready binaries see releases page.

## Using as library for your developments

You can use this repository as a library in your Go projects.
To do this,
you need to import the package **github.com/kt-soft-dev/kt-cli/pkg** and use the functions provided by the client.


## Making API request

To make an API request, you can use -act.method flag to specify the method of the request. For example: 

```bash
ktcloud -act.method=test.test
```

Output will be like this:

```bash
2024/04/01 06:37:17 {"ok":true}
```

Parameters can be passed using **-params** flag.
Value should be a string with space-separated key-value pairs.

For example:

```bash
ktcloud -act.method=test.test -params="param1=value1 param2=value2"'
```

In this example params are just stubs and will be ignored. To get known about parameters for specific method, please read the API documentation.

## Output modes

Output can be displayed in different modes. By default, output is displayed in usual **log.Println** format like this:

```bash
2024/04/01 06:37:17 {"ok":true}
```

**-output** flag can be used to specify output mode. Currently supported modes are:

- **0** - log with timestamp
- **1** - plain log (simple output, no timestamp)
- **2** - just like plain log but without new line at the end

## Contributing

Contributions are welcome. Please feel free to submit a pull request or open an issue on the GitHub repository.

## License

See the `LICENSE` file.