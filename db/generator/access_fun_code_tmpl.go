package generator

const findCodeYamlTemplate = `steps:
    - lookup:
        nout: stmt
        obj: db
        name: preparedCache
        key: {val: "{{ .Name }}"}
    - call:
        nout: [values, err]
        func: "{{ .Name }}ParseParams"
        args: [reqeust]
        err_returns: [nil, err]
    - call:
        nout: [rows, err]
        obj: stmt
        func: Query
        args: ["values..."]
        err_returns: [nil, err]
        clean: 
            obj: rows
            func: Close
    - var:
        name: results
        type: "[]{{.ModelName}}"
    - repeat_cond:
        cond: {call: {func: Next, obj: rows}}
        body:
        - var:
            name: item
            type: "{{.ModelName}}"
        - call:
            nout: err
            obj: rows
            func: Scan
            args: [{{ .ScanAttributes }}]
            err_returns: [nil, err]
        - call:
            out: results
            func: append
            args: [results, item]
    - return: results, nil`

const updateCodeYamlTemplate = `steps:
    - lookup:
        nout: stmt
        obj: db
        name: preparedCache
        key: {val: "{{ .Name }}"}
    - call:
        nout: [values, err]
        func: {{ .Name }}ParseParams
        args: request
        err_returns: [0, err]
    - call:
        nout: [result, err]
        obj: stmt
        func: Exec
        args: "values..."
        err_returns: [0, err]
    - call:
        nout: [rowsAffected, err]
        obj: result
        func: RowsAffected
    - return: [rowsAffected, err]`
const addCodeYamlTemplate = `steps:
    - lookup:
        nout: stmt
        obj: db
        name: preparedCache
        key: {val: "{{ .Name }}"}
    - call:
        nout: [values, err]
        func: {{ .Name }}ParseParams
        args: request
        err_returns: [0, err]
    - var:
        name: id
        type: int64
    - call:
        nout: [result, err]
        obj: stmt.QueryRow()
        func: Exec
        args: ["&id"]
    - return: [id, err]`
const addOrReplaceCodeYamlTemplate = `steps:
    - lookup:
        nout: stmt
        obj: db
        name: preparedCache
        key: {val: "{{ .Name }}"}
    - call:
        nout: [values, err]
        func: {{ .Name }}ParseParams
        args: request
        err_returns: [0, err]
    - var:
        name: id
        type: int64
    - var:
        name: inserted
        type: bool
    - call:
      nout: [err]
      obj: stmt.QueryRow()
      func: Exec
      args: ["&id", "&inserted"]
    - return: [id, inserted, err]`
const deleteCodeYamlTemplate = `steps:
    - lookup:
        nout: stmt
        obj: db
        name: preparedCache
        key: {val: "{{ .Name }}"}
    - call:
        nout: [values, err]
        func: {{ .Name }}ParseParams
        args: request
        err_returns: [0, err]
    - call:
        nout: [result, err]
        obj: stmt
        func: Exec
        args: "values..."
        err_returns: [0, err]
    - call:
        nout: [rowsAffected, err]
        obj: result
        func: RowsAffected
    - return: [rowsAffected, err]`
