import { ArrowUpIcon } from "@heroicons/react/solid";
import { SwitchHorizontalIcon } from "@heroicons/react/outline";
import GitHubIcon from '@mui/icons-material/GitHub';
import LinkedInIcon from '@mui/icons-material/LinkedIn';
import EmailOutlinedIcon from '@mui/icons-material/EmailOutlined';
import React from "react";
import { Link } from "react-scroll";

function Footer() {
    return (
        <footer className="bg-gray-800 md:sticky top-0 z-10">
            <div className="container mx-auto flex flex-wrap p-5 flex-col md:flex-row items-center">
                <div className="title-font font-medium text-white mb-4 md:mb-0">
                    <span className="ml-3 mr-4 text-xl">Connect</span>
                    <SwitchHorizontalIcon className="inline-flex w-8 h-5 text-gray-500 mb-1" />
                </div>
                <nav className="md:mr-auto md:ml-4 md:py-1 md:pl-4 md:border-l md:border-gray-700 flex flex-wrap items-center text-base justify-center">
                    <a href="https://github.com/rfabreu" target="_blank" className="mr-5 hover:text-white" rel="noreferrer">
                        <GitHubIcon />
                    </a>
                    <a href="https://www.linkedin.com/in/rafael-a-gomes/" target="_blank" className="mr-5 hover:text-white" rel="noreferrer">
                        <LinkedInIcon />
                    </a>
                    <a href="mailto:raabreugomes@gmail.com" className="mr-5 hover:text-white" rel="noreferrer">
                        <EmailOutlinedIcon />
                    </a>
                </nav>
                <Link to="about" activeClass="active" spy={true} smooth={true} offset={-70} duration={2000} style={{cursor:"pointer"}} className="inline-flex items-center bg-gray-800 border-0 py-1 px-3 focus:outline-none hover:bg-gray-700 rounded text-base mt-4 md:mt-0">
                    Home
                    <ArrowUpIcon className="w-4 h-4 ml-1" />
                </Link>
            </div>
        </footer>
    );
}

export default Footer;