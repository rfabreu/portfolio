import { ArrowRightIcon, TerminalIcon } from "@heroicons/react/solid";
import React from "react";
import Resume from "../docs/Rafael_Abreu_Resume.pdf";
import { Link } from "react-scroll";

function Navbar() {
    return (
        <header className="bg-gray-800 sticky top-0 z-10">
            <div className="container mx-auto flex flex-wrap p-5 flex-col md:flex-row items-center">
                <a className="title-font font-medium text-white mb-4 md:mb-0">
                    <Link to="about" activeClass="active" spy={true} smooth={true} offset={-70} duration={2000} style={{ cursor: "pointer" }} className="ml-3 text-xl">
                        <TerminalIcon className="inline-flex w-8 text-gray-500 mb-1" />
                        Rafael Abreu
                    </Link>
                </a>
                <nav className="md:mr-auto md:ml-4 md:py-1 md:pl-4 md:border-l md:border-gray-700 flex flex-wrap items-center text-base justify-center">
                    <Link to="projects" activeClass="active" spy={true} smooth={true} offset={-70} duration={1000} style={{ cursor: "pointer" }} className="mr-5 hover:text-white">
                        Projects
                    </Link>
                    <Link to="skills" activeClass="active" spy={true} smooth={true} offset={-70} duration={1000} style={{ cursor: "pointer" }} className="mr-5 hover:text-white">
                        Skills
                    </Link>
                    <a href={Resume} download="Rafael_Abreu_Resume.pdf" className="mr-5 hover:text-white">
                        Resume
                    </a>
                    {/* <a href="#testimonials" className="mr-5 hover:text-white">
                        Testimonials
                    </a> */}
                </nav>
                <Link to="contact" activeClass="active" spy={true} smooth={true} offset={-70} duration={2000} style={{ cursor: "pointer" }} className="inline-flex items-center bg-gray-800 border-0 py-1 px-3 focus:outline-none hover:bg-gray-700 rounded text-base mt-4 md:mt-0">
                    Contact
                    <ArrowRightIcon className="w-4 h-4 ml-1" />
                </Link>
            </div>
        </header>
    );
}

export default Navbar;