import React, {useState, useEffect} from 'react';
import './index.css';

const Sleep = ms => new Promise(
    resolve => setTimeout(resolve, ms)
);


const BubbleSort = () => {
    const [arr, setArr] = useState([]);
    const [stop, setStop] = useState(0);
    const [pause,setPause] = useState(true);
    let speed = 10, dataRange = 65, dataLength = 80
    window.localStorage.setItem('pause',pause)
    window.localStorage.getItem('speed') === null ? window.localStorage.setItem('speed', speed) : speed = window.localStorage.getItem('speed')
    window.localStorage.getItem('dataRange') === null ? window.localStorage.setItem('dataRange', dataRange) : dataRange = window.localStorage.getItem('dataRange')
    window.localStorage.getItem('dataLength') === null ? window.localStorage.setItem('dataLength', dataLength) : dataLength = window.localStorage.getItem('dataLength')

    const ResetMainArray = () => {
        setArr([])
        for (let i = 0; i < dataLength; i++) {
            setArr(arr => [...arr, Math.floor(Math.random() * dataRange)+1])   
        }
        for (let i = 0; i < arr.length; i++) {
            document.getElementById(i).style.backgroundColor = 'white'
            document.getElementById(i).style.height = `${arr[i] * 10}px`;
        }
    }

    const SliderFunction = (e,id) => {
        document.getElementById(id).innerHTML = e.target.value
        window.localStorage.setItem(`${id}`, e.target.value)
        if (id !== "speed") {
            ResetMainArray();
            window.localStorage.setItem('stop',1)
            setStop(1)
        }
    }

    const Sorting = async () => {
        window.localStorage.setItem('stop',0)
        let array = arr
        for (let i = 0; i < array.length; i++) {
            for (let j = 0; j < array.length-i-1; j++) {
                document.getElementById(j).style.height = `${array[j] * 10}px`; 
                if (array[j] > array[j+1]) {
                    let temp = array[j]
                    array[j] = array[j+1]
                    array[j+1] = temp
                } 
                console.log(window.localStorage.getItem('stop'));
                await Sleep(window.localStorage.getItem('speed'))
                if (window.localStorage.getItem('stop') === 1 ) {
                    return
                }
                await Sleep(window.localStorage.getItem('stop'))
                document.getElementById(j).style.backgroundColor = 'red'
                document.getElementById(j+1).style.backgroundColor = 'yellow'   
                document.getElementById(j).style.height = `${array[j] * 10}px`; 
            }
            for (let k = 0; k < array.length; k++) {
                if (k <= array.length - i - 2) {
                    document.getElementById(k).style.backgroundColor = 'white'
                } else {
                    document.getElementById(k).style.backgroundColor = 'yellow'          
                }
                document.getElementById(k).style.height = `${array[k] * 10}px`;
            }
        }   
        alert('Bubble Sort!')
    }

    useEffect(()=> {
        ResetMainArray()
    }, [])

    return (
        <div className="Bubble">
            <div className="OptionsPanel">
                <div className="slidecontainer">
                    <h2>Speed</h2>
                    <input onChange={(e) => SliderFunction(e, 'speed')}  id="myRangeSpeed" type="range" min="1" max="1000" defaultValue={speed} className="slider"/>
                    <p>Value: <span id="speed">{speed}</span></p>
                </div>
                <div className="slidecontainer">
                    <h2>Data range</h2>
                    <input onChange={(e) => SliderFunction(e, "dataRange")}  id="myRangeData" type="range" min="0" max="75" defaultValue={dataRange} className="slider"/>
                    <p>Value: <span id="dataRange">{dataRange}</span></p>
                </div>
                <div className="slidecontainer">
                    <h2>Data length</h2>
                    <input onChange={(e) => SliderFunction(e, "dataLength")}  id="myRangeLength" type="range" min="0" max="85" defaultValue={dataLength} className="slider"/>
                    <p>Value: <span id="dataLength">{dataLength}</span></p>
                </div>
                <div className="slidecontainer">
                    <button onClick={()=> ResetMainArray()}>Reset</button>
                    {/* <button onClick={(e) =>{     
                        // if (pause) {
                        //     window.localStorage.setItem('speed', 100000)
                        //     e.target.innerHTML = "Resume"
                        // } else {
                        //     window.localStorage.setItem('speed', 0)
                        //     e.target.innerHTML = "Pause"
                        // }
                        // setPause(!pause)
                        }}>Pause</button> */}
                </div>
                <div className='qwe'>
                    <button className="mainButton" onClick={() => Sorting()}>Sort It</button>
                </div>

            </div>
            <div className="Algorithm">
                {arr?.map((item, key) => {
                    return <div key={key} id={key} style={{height:`${item*10}px`}} className="box"></div>  
                })}
            </div>
        </div>
    );
}

export default BubbleSort;
