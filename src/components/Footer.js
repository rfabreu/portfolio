import { ArrowUpIcon } from "@heroicons/react/solid";
import { SwitchHorizontalIcon } from "@heroicons/react/outline";
import GitHubIcon from '@mui/icons-material/GitHub';
import LinkedInIcon from '@mui/icons-material/LinkedIn';
import React from "react";

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
                    {/* Stackoverflow Icon */}
                    <a href="https://stackoverflow.com/users/18742772/rafael" target="_blank" className="mr-5 hover:text-white" rel="noreferrer">
                    <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 32 32" width="24" height="24" class="MuiSvgIcon-root MuiSvgIcon-fontSizeMedium css-i4bv87-MuiSvgIcon-root" focusable="false" aria-hidden="true"><path d="M28.16 32H2.475V20.58H5.32v8.575h19.956V20.58h2.884z"></path><path d="M8.477 19.8l13.993 2.923.585-2.806-13.993-2.923zm1.832-6.704l12.94 6.04 1.208-2.572-12.94-6.08zm3.586-6.353l10.99 9.12 1.832-2.183-10.99-9.12zM20.99 0l-2.3 1.715 8.536 11.46 2.3-1.715zM8.166 26.27H22.43v-2.845H8.166v2.845z"></path></svg>
                    </a>
                </nav>
                <a href="#" className="inline-flex items-center bg-gray-800 border-0 py-1 px-3 focus:outline-none hover:bg-gray-700 rounded text-base mt-4 md:mt-0">
                    Home
                    <ArrowUpIcon className="w-4 h-4 ml-1" />
                </a>
            </div>
        </footer>
    );
}

export default Footer;