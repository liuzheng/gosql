# gosql

A SQL generate, like django ORM.
GoSQL will not support DB table `ForeignKey`, `ManyToManyField` and other table 

## Install

```shell
$ go get github.com/liuzheng/gosql
```

## Usage

In your go code you need add some tag in the struct, like:
```go
type User struct {
	id         uint16    `mysql:"SMALLINT,NOT_NULL,AUTO_INCREMENT,PRIMARY_KEY"`
	Name       string    `mysql:"varchar(128)"`
	Avatar     string    `mysql:"varchar(128)"`
	password   string    `mysql:"varchar(128)"`
	CreateTime time.Time `mysql:"timestamp"`
}
```
Then run this command `gosql -makemigrations .` will generate the SQL file under migrations folder in the local path.

## Road map

In the 1.x version, I'm plan only support one tag, multi sql tag will be ignored, and generate the SQL migration file.

In the 2.x version, I plan to make it support only one tag named `gosql` and generate SQL file with command flag `-sqltype {mysql,sqlite}`.
With the tag `gosql` and `mysql`(for e.g.) defined, gosql will only support `gosql` tag.

