import React from "react";
import { Link } from "react-router-dom";

function NavBar() {
    return(
    <>
        <nav className="navbar navbar-expand-lg navbar-light bg-light">
            <div className="container-fluid">
                <Link to ="/" style={{textDecoration: 'none'}}>
                    <a className="navbar-brand" href="#">Navbar</a>
                </Link>
                <button className="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#navbarNav" aria-controls="navbarNav" aria-expanded="false" aria-label="Toggle navigation">
                    <span className="navbar-toggler-icon"></span>
                </button>
                <div className="collapse navbar-collapse" id="navbarNav">
                    <ul className="navbar-nav">
                        <li className="nav-item">
                            <Link to="/" style={{textDecoration: 'none'}}>
                                <a className="nav-link active" aria-current="page" href="#">Home</a>
                            </Link>
                        </li>
                        <li className="nav-item">
                            <Link to="/reportes" style={{textDecoration: 'none'}}>
                                <a className="nav-link" href="#">Reportes</a>
                            </Link>
                        </li>
                    </ul>
                </div>
            </div>
        </nav>
    </>
    )
}

export default NavBar;