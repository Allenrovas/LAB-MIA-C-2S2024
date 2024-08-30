import React from "react";
import NavBar from "../components/NavBar";
import Consola from "../components/Editor";
import Service from "../Services/Service";

function Index(){
    const [value, setValue] = React.useState("");
    const [response, setResponse] = React.useState("");

    const changeText = (text) => {
        setValue(text);
    }

    const handlerClick =  () => {
        console.log(value);
        if (value === ""){
            alert("No puedes enviar un comando vacío");
            return;
        } 
        Service.analisis(value)
        .then((res) => {
            setResponse(res.data);
        })
        .catch((err) => {
            console.error(err);
        });
    }

    const handlerLimpiar = () => {
        if (value === ""){
            alert("No puedes limpiar un campo vacío");
            return;
        }
        changeText("");
        setResponse("");
    }

    const handleLoadClick = () => {
        const input = document.createElement('input');
        input.type = 'file';
        input.addEventListener('change', handleFileChange);
        input.click();
    }

    const handleFileChange = (e) => {
        const file = e.target.files[0];
        const reader = new FileReader();
        reader.onload = (e) => {
            const text = e.target.result;
            changeText(text);
        }
        reader.readAsText(file);
    }

    return (
        <>
            <NavBar />
            <h1>Proyecto 1 - MIA - 202004745</h1>
            <Consola text={"Consola de Entrada"} handlerChange={changeText} value={value}/>
            <div class="container">
                <button type="button" class="btn btn-primary" onClick={handlerClick}>Enviar</button>
                <button type="button" class="btn btn-danger" onClick={handlerLimpiar}>Limpiar</button>
                <button type="button" class="btn btn-success" onClick={handleLoadClick}>Cargar Archivo</button>
            </div>
            <Consola text={"Consola de Salida"} handlerChange={changeText} value={response} readOnly={true}/>
        </>
    )
}

export default Index;