# Linguagem ASM
#### Por: Arthur Andrade
Este documento especifica detalhes de uso da linguagem ASM(Esse nome foi gerado por uma IA geradora de nomes)

## Compilar
Para rodar o arquivo ASM, basta fazer a build do executavel com `go build -o asm cmd/main.go`.
Após ter o executável, escreva seu progrmaa ASM e chame o executabel passando o nome do arquivo como primeiro argumento.
Exemplo:
```
./asm program.asm 
```

O outputa sairá no mesmo diretório do input, mas você pode ainda especificar o ponto de saida do arquivo como no exemplo abaixo.
```
./asm program.asm program.bin
```

## Números
Literais de números podem ser escritos de quatro formas diferentes desde que sigam seus prefixos. Sendo elas binários, decimais, octais e hexadecimais.
Segue o exemplo abaixo do número onze:
```
0b1011 // 11 em número binario, prefixo '0b'
0o13   // 11 em número octal, prefixo '0o'
11     // 11 em número decimal, sem prefixo
0xb    // 11 em número hexadecimal, prefixo '0x'
```

## Endereços de memória
Endereços de memória são simples e basta colocar o prefixo `m`, seguido por números inteiros em formato decimal
para apontar para a célula requisitada. Segue o exemplo abaixo:
```
m4  // Aponta para o 4º endereço de memória.
m12 // Aponta para o 12º endereço de memória.
```

## Labels
As labels marcam pontos onde futuramente minmônicos de comparação poderão apontar como ponto de 'jump' no código.
Para criar um label, basta escrevê-la com o prefixo `#` e colocar dois pontos após seu final. Nomes de labels só
podem conter letras minusculas, maisculas, numeros e underline(_). O primeiro caractere de uma label n pode ser
um numero nem um underline. Segue a criação de quatro labels abaixo:
```
#here:         
#anotherLabel: 
#here2:
#he_re3:
```
# Mnemônicos da Arquitetura

## Memory manipulations

### 1. `GET` - **GET**
- Copia um valor para o acumulador, este valor pode ser um literal ou um endereço de memória.
- **Format**: `GET <MEMORY_ADDRESS | NUMBER_LITERAL>`
- **Exemplo de uso**:

  ```asm
  GET 3      // Puts a 3 into accumulator.
  GET m3     // Gets value from 3º byte to accumulator.

### 2. `SET` - **SET**
- Salva o valor do acumulador dentro de um endereço de memória.
- **Format**: `SET <MEMORY_ADDRESS>`
- **Exemplo de uso**:

  ```asm
  SET m1      // Writes value from accumulator to first byte.
  SET m2      // Writes value from accumulator to second byte.

### 3. `CPY` - **COPY**
- Copia o valor de um acumulador ou literal para outro acumulador ou memória.
- **Format**: `CPY <MEMORY_ADDRESS> <MEMORY_ADDRESS | NUMBER_LITERAL>`
- **Exemplo de uso**:

  ```asm
  CPY m3, 3   // Puts a 3 into 3º byte.
  CPY m3, m4  // Copies from 4º byte to 3º byte

## Simple operations

### 4. `INC` - **Increment**
- Incrementa o valor do acumulador.
- **Exemplo de uso**:

  ```asm
  INC   // Incrementa um do acumulador

### 5. `DEC` - **Decrement**
- Decrementa o valor do acumulador.
- **Exemplo de uso**:

  ```asm
  DEC   // Decrementa um do acumulador

### 6. `NEG` - **Negate**
- Inverte o sinal do valor no acumulador.
- **Exemplo de uso**:

  ```asm
  NEG   // Inverte o sinal do acumulador

### 7. `NOT` - **NOT**
- Aplica a operação 'not' nos bits do acumulador.
- **Exemplo de uso**:

  ```asm
  NOT   // Aplica o 'not' nos bits do acumulador

## Operations

### 8. `ADD` - **ADD**
- Soma o valor do acumulador com o valor de um endereço de memória.
- **Format**: `ADD <MEMORY_ADDRESS>`
- **Exemplo de uso**:

  ```asm
  ADD m1 // Adiciona o conteúdo de m1 no acumulador

### 10. `AND` - **AND**
- Aplica a operação lógica AND entre o acumulador e o valor de um endereço de memória.
- **Format**: `AND <MEMORY_ADDRESS>`
- **Exemplo de uso**:

  ```asm
  AND m1 // Aplica a operação 'and' entre o conteúdo de m1 e o acumulador. Salva no acumulador

### 11. `OR` - **OR**
- Aplica a operação lógica OR entre o acumulador e o valor de um endereço de memória.
- **Format**: `OR <MEMORY_ADDRESS>`
- **Exemplo de uso**:

  ```asm
  OR m1 // Aplica a operação 'or' entre o conteudo de m1 e o acumulador. Salva no acumulador

### 12. `XOR` - **XOR**
- Aplica a operação lógica XOR entre o acumulador e o valor de um endereço de memória.
- **Format**: `XOR <MEMORY_ADDRESS>`
- **Exemplo de uso**:

  ```asm
  XOR m1 // Aplica a operação 'xor' entre o conteúdo de m1 e o acumulador. Salva no acumulador

## Loops and comparations

### 13. `SUB` - **Compare**
- Aplica subtração entre o acumulador e o valor da memória passada.
- **Exemplo de uso**:

  ```asm
  SUB m1 // Aplica a operação 'xor' entre o conteudo de m1 e o acumulador. Salva no acumulador

### 14. `JMP` - **Jump**
- Faz o salto para a 'label' especificado.
- **Format**: `JMP <LABEL>`
    - **Exemplo de uso**:

  ```asm
  #here:
    JMP #here // Cria um loop infinito aonde sempre volta para 'here'.

### 15. `JIZ` - **Jump if zero**
- Faz o salto para a 'label' especificado se o acumulador for zero.
- **Format**: `JIZ <LABEL>`
- **Exemplo de uso**:

  ```asm
  #here:
    JIZ #here // Cria um loop infinito se houver um zero no acumulador.

### 16. `JIN` - **Jump if negative**
- Faz o salto para a 'label' especificado se o acumulador for negativo.
- **Format**: `JIN <LABEL>`
- **Exemplo de uso**:

  ```asm
  #here:
    JIN #here // Cria um loop infinito se houver um valor negativo no acumulador.

## Runtime actions
### 19. `HLT` - **Halt**
- Interrompe a execução do programa.
- **Exemplo de uso**:

  ```asm
  HLT  // Interrompe a execução
