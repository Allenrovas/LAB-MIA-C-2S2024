import React from "react";
import Editor from "@monaco-editor/react";
import "../styles/Editor.css";

function Consola(props){
    const handlerChangeEditor = (newValue, e) => {
        props.setCodigo(newValue);
    }

    return(
        <>
            <div className="mb-3" editor-react>
                <label
                htmlFor="exampleFormControlTextarea1"
                style={{ fontSize: "30px"}}
                >
                {props.text}
                </label>
                <Editor
                height="50vh"
                theme="vs-dark"
                defaultLanguage="typescript"
                value={props.value}
                onChange={handlerChangeEditor}
                options={{
                    selectOnLineNumbers: true,
                    automaticLayout: true,
                    lineNumbers: "on",
                    readOnly: props.readOnly,
                }}
                />
            </div>
        </>

    );
}

export default Consola;