# randstring

Generate random strings based on charset and composition requirements.

```shell
go get -u github.com/absurdlab/pkg/randstring
```

## Features
- Precise control on output character composition
- Baked in character sets

## Usage

```go
// Generate a string that has 5 lower case, 2 upper case, 
// 1 number and 1 special characters.
password := randstring.Generate(
    randstring.Charsets.LowerCaseAlpha(5),
    randstring.Charsets.UpperCaseAlpha(2),
    randstring.Charsets.Numeric(2),
    randstring.Charsets.Special(1)
)
```