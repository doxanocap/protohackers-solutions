import { useNavigate } from 'react-router-dom';
import React from "react";
import "./index.css"

const Header = () => {
    let menu;
    //const navigate = useNavigate();
    // const logOut = async (e) => {
        // e.preventDefault();
        // const response = await fetch("http://localhost:8080/api/logout", {
        //     headers: { 'Content-Type': 'application/json' },
        //     credentials: 'include',
        // });
        // const data = await response.json()
        // if (data.message === "deleted cookie") {
        //     console.log("blablabla logout")
        //     navigate("/")
        // } else {
        //     console.log("щщс не работает")
        // }
        // console.log(data)
        // window.location.reload();
    // }
    // if (name === "" || typeof (name) === "undefined") {
    //     menu = (
    //         <nav>

    //             <a href="/posts" onClick={() => { navigate('/posts') }}>Posts</a>
    //             <a href="/login" onClick={() => { navigate('/login') }}>Login</a>
    //             <a href="/register" onClick={() => { navigate('/register') }}>Register</a>
    //         </nav>
    //     )
    // } else {
    //     menu = (
    //         <nav>

    //             <a href="/posts" onClick={() => { navigate('/posts') }}>Posts</a>
    //             <a href="/" onClick={logOut}>Logout</a>
    //             <a href='/account-preferences'>{name} {surname}</a>
    //         </nav>
    //     )
    // }
    return (
        <div className="pagesHeader">
            <h2 onClick={() => console.log('/')} id="mainp">CV-review</h2>
            <nav>
                <a href="/posts" onClick={() => { console.log('/posts') }}>Posts</a>
                <a href="/login" onClick={() => { console.log('/login') }}>Login</a>
                <a href="/register" onClick={() => { console.log('/register') }}>Register</a>
            </nav>
        </div>
    )
};

export default Header;