## Yusra
Yusra is a CLI tool that helps you navigate your JSON data

### Motivation
JSON is an incredibly popular encoding format, widely used by almost every API today. However, JSON can suck

- It's large because it's text-based.
- It's heavy to process (have you ever tried opening a 10MB JSON file?).
- You can't query it unless it's loaded into a database like MongoDB.

To address these issues, I created Yusra. With Yusra, you can load a JSON file and query it using SQL. This allows me to gain valuable insights into the dataset and makes working with it much easier. Additionally, I can perform joins by loading multiple sets, provided there is a foreign key in the JSON dataset.

### Installation
**__You need to have [Golang](https://go.dev/doc/install) and [cmake](https://cmake.org/download/) installed locally__**

You can the [releases](https://github.com/mohammedbalila/jsondb/releases/) to download the binary directly or clone the repo and build it locally

```bash
git clone https://github.com/mohammedbalila/yusra.git && cd yusra 
```

```bash
make build
```

```bash
chmod +x yusra && mv /usr/local/bin/
```

### How to use
```bash
yusra load filename.json
# this will load the file into yusra's database and gives you a prompt that you can run SQL queries on
```
> Please note that the JSON file has to be an array of objects, this is a limitation yusra will overcome very soon!
```bash
yusra>select * from filename;
```

```bash
yusra>help
# this will show the help message
```

| Command               | Description                           |
|-----------------------|---------------------------------------|
| `exit`                | To exit                               |
| `version`             | To show version                       |
| `load file.json`      | To load a new JSON file               |
| `files`               | To list loaded set                    |
| `info dataset_name`   | To get info about a loaded dataset    |


### How it works
Yusra is smart, it knows that in order to defeat JSON she has to strip it from its power (text-based). The JSON file is loaded into a SQLite database
and then exposed for the end user to work with.

### How to contribute
Please refer to the [Code of Conduct](./CODE_OF_CONDUCT). 

If you encounter any issues, please create a new issue with a brief description and the steps to reproduce the problem. Including environment details such as your OS and Golang version would be very helpful.

The code is still a work in progress, so pull requests are always welcome!