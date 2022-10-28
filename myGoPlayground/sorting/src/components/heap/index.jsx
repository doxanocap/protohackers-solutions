import React, {useState, useEffect} from 'react';
import './index.css';

const Sleep = ms => new Promise(
    resolve => setTimeout(resolve, ms)
);


const HeapAlgorithm = () => {

    return (
        <div className="Heap">
            <div className="Floor">
                <div className="Nodes">1</div>
            </div>
            <div className="Floor">
                <div className="Nodes">2</div>
                <div className="Nodes">3</div>
            </div>
            <div className="Floor">
                <div className="Nodes">4</div>
                <div className="Nodes">5</div>
                <div className="Nodes">6</div>
                <div className="Nodes">7</div>
            </div>
            <div className="Floor">
                <div className="Nodes">8</div>
                <div className="Nodes">9</div>
                <div className="Nodes">10</div>
                <div className="Nodes">11</div>
                <div className="Nodes">12</div>
            </div>
        </div>
    );
}

export default HeapAlgorithm;
