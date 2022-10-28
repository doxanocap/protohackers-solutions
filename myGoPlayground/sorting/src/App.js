import { BrowserRouter, Route, Routes } from 'react-router-dom';
import { React, useState, useEffect } from "react";
import BubbleSort from './components/bubbleSort'
import HeapAlgorithm from './components/heap' 
import Header from './components/header';
import './App.css';

const App = () => {
  return (
    <BrowserRouter>
    <Header />
    <Routes>
      <Route path={"/bubble-sort"} element={<BubbleSort />} />
      <Route path={"/"} element={<HeapAlgorithm />} />
    </Routes>
  </BrowserRouter>
  );

}

export default App;
