# P1 - final


### Compilar
```
$ make all
```


### Testar na raiz do projeto
```
$ exp programa.lpn output.asm
$ asmp1 output.asm output.bin
$ neander output.bin
```


## Programa LPN

Esta e a sintaxe basica para funcionar.
Os operadores permitidos `+`, `-`, `&` e `|`, representando soma, subtração, and e or respetivamente.

A unica funcao permitida é a main, que é criada como demonstrado abaixo, o print coloca o valor resultado na posicao de memoria 200 do Neander

```text
func main() {

    def foo = 1
    foo += 1

    def bar = 1 + 3
    bar = 3

    def chaos = foo + bar

    print(foo)
}
```
