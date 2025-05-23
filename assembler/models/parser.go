package models

import (
    "ASM/neanderExecutor"
)

type (
    MemHeap struct {
        content map[uint16]int16
        last    uint16
    }

    Parser struct {
        Filename string
        labels   map[string]uint16
        memHep   MemHeap
        tokens   []Token
        output   Ast
        cursor   int
    }
)

func NewParser(filename string, tokens []Token) Parser {
    // Pega o index da ultima constante declarada
    l := GetLastConstant()
    constants := GetBuiltinConstants()
    return Parser{
        Filename: filename,
        labels:   map[string]uint16{},
        memHep: MemHeap{
            content: constants,
            last:    l,
        },
        output: Ast{},
        tokens: tokens,
        cursor: 0,
    }
}

func (this *Parser) AllocNum(num int16) uint16 {
    this.memHep.last++
    where := this.memHep.last
    this.memHep.content[where] = num
    return where - NeanderPadding + JmpConstantsSize
}

func (this *Parser) WriteProgram() ([]uint8, error) {

    var vec [][]uint16               // Guarda os bytes de cada statement
    var stmtSizes int                // Tamanho em tempo real do programa
    resolveLabel := map[int]string{} // Labels que devem ter posições recalculadas após a definição do header
    var reviewOffset []uint16        // Posicoes que devem ter posições recalculadas após a definição do header
    for _, stmt := range this.output.Statements {

        // Se for um statement de label, adiciona nos bytecode labels
        if stmt.GetTitle() == "LabelDeclStmt" {
            labelStmt := stmt.(LabelDeclStmt)
            this.labels[labelStmt.LabelName] = uint16(stmtSizes)
        }

        // Transforma o statement em bytecode
        bytes, mems, err := stmt.WriteMemASM()
        if err != nil {
            return nil, err
        }

        // Se for um jump, marca a label como pendencia de recalculo de posição
        if stmt.GetTitle() == "JumpStmt" {
            jmpStmt := stmt.(JumpStmt)
            resolveLabel[stmtSizes+len(bytes)-1] = jmpStmt.TargetLabelName
        }

        // Corrige as posições relativas para o programa inteiro, não o stmt
        for i := range mems {
            mems[i] += uint16(stmtSizes) - 1
        }

        // Adiciona no array de posições de memórias para ser recalculadas com o offset
        reviewOffset = append(reviewOffset, mems...)

        stmtSizes += len(bytes)
        vec = append(vec, bytes)
    }

    // Gera o array real do programa
    var program []uint16
    for _, bytes := range vec {
        program = append(program, bytes...)
    }

    constants := this.memHep.content // Constantes do programa
    constantsCount := uint16(len(constants))
    neanderPrefix := []uint16{1, 1} // Prefixo de 2 bytes de todos os programas neander
    PaddingSize := uint16(len(neanderPrefix)) + constantsCount

    // Recalcula a posicao das labels
    for k := range this.labels {
        if this.labels[k] != 0 {
            this.labels[k] += PaddingSize
        }
    }

    // Aplica as novas posicoes das labels nos lugares onde elas foram chamadas
    for k, v := range resolveLabel {
        if this.labels[v] != 0 {
            program[k] = this.labels[v]
            continue
        }
        return nil, GetUnkownLabelErr(this.Filename, v)
    }

    // Recalcula as posições de memória com o offset
    for _, mem := range reviewOffset {
        program[mem] += PaddingSize
    }

    // Garante q o programa não vai tentar executar o espaço reservado para constantes
    neanderPrefix = append(neanderPrefix, neander.JMP, PaddingSize)
    // Adiciona as constantes e os espaços de memória reservados.
    neanderPrefix = append(neanderPrefix, make([]uint16, constantsCount)...)
    for k, v := range constants {
        neanderPrefix[k] = uint16(uint8(v))
    }

    // Transforma tudo em uint16
    neanderPrefix = append(neanderPrefix, program...)
    final := make([]uint8, len(neanderPrefix)*2)
    for i, num := range neanderPrefix {
        final[i*2+1] = uint8(num >> 8)
        final[i*2] = uint8(num)
    }

    // Marca o fim do programa
    endAt := len(final)

    // Itera pra ter 516 de tamanho
    if endAt < 516 {
        final = append(final, make([]uint8, 516-endAt)...)
    }

    return final, nil
}

func (this *Parser) Inject(stmts ...Stmt) {
    this.output.Statements = append(this.output.Statements, stmts...)
}

func (this *Parser) Inspect() {
    this.output.Inspect()
}

func (this *Parser) Get(n int) *Token {
    if this.cursor+n >= len(this.tokens) {
        return nil
    }
    return &this.tokens[this.cursor+n]
}

func (this *Parser) Consume(n int) {
    if this.cursor >= len(this.tokens) {
        return
    }
    this.cursor += n
}

const (
    NoSpaceMode = iota
    OptionalSpaceMode
    MandatorySpaceMode
)

func (this *Parser) HasNextConsume(spaceMode int, kinds ...TokenKindEnum) *Token {
    if spaceMode < NoSpaceMode || spaceMode > MandatorySpaceMode {
        panic("invalid argument in function 'HasNextConsume'")
    }
    for findSpace := false; ; {
        token := this.Get(0)
        if token == nil {
            // Fim dos tokens sem encontrar um tipo esperado
            return nil
        }

        for _, kind := range kinds {
            if token.Kind == kind {
                // Se espaços eram obrigatórios mas não foram encontrados, falha
                if spaceMode == MandatorySpaceMode && !findSpace {
                    return nil
                }
                this.Consume(1)
                return token
            }
        }

        if token.Kind == TOKEN_SPACE {
            findSpace = true
            this.Consume(1)
            continue
        }

        // Se espaços não eram permitidos ou eram obrigatórios e encontrou outro token, falha
        if spaceMode == NoSpaceMode || spaceMode == MandatorySpaceMode {
            return nil
        }

        return nil // Qualquer outro caso não esperado falha
    }
}
