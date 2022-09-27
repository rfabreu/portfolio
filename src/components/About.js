import React from "react";
import hero from "../images/hero.png";
import { Link } from "react-scroll";

function About() {
    return (
        <section id="about">
            <div className="container mx-auto flex px-10 py-20 md:flex-row flex-col items-center">
                <div className="lg:flex-grow md:w-1/2 lg:pr-24 md:pr-16 flex flex-col md:items-start md:text-left mb-16 md:mb-0 items-center text-center">
                    <h1 className="title-font sm:text-4xl text-3xl mb-4 font-medium text-white">
                        Hi, I'm Rafael.
                        <br className="hidden lg:inline-block" />
                        <span className="title-font sm:text-4xl text-3xl mb-4 font-medium text-white">Software Engineer and Full-stack Web Developer.</span>
                    </h1>
                    <p className="mb-8 leading-relaxed">
                    I studied at the University of Toronto coding Bootcamp, and recently acquired a certificate from Red Had in OpenShift and Kubernetes. Throughout the years I've demonstrated success in the technology, television, and education industries. Where through teamwork and leadership I contributed to implementing innovative solutions. My goal is to help you build great things through good Software Engineering practices.
                    </p>
                    <div className="flex justify-center">
                        <Link to="contact" activeClass="active" spy={true} smooth={true} offset={-70} duration={2000} style={{cursor:"pointer"}} className="inline-flex text-white bg-green-500 border-0 py-2 px-6 focus:outline-none hover:bg-green-600 rounded text-lg">
                            Work With Me
                        </Link>
                        <Link to="projects" activeClass="active" spy={true} smooth={true} offset={-70} duration={1000} style={{cursor:"pointer"}} className="ml-4 inline-flex text-gray-400 bg-gray-800 border-0 py-2 px-6 focus:outline-none hover:bg-gray-700 hover:text-white ronded text-lg">
                            See My Work
                        </Link>
                    </div>
                </div>
                <div className="lg:max-w-lg lg:w-full md:w-1/2 w-5/6">
                    <img className="object-cover object-center rounded" alt="hero" src={hero} />
                </div>
            </div>
        </section>
    );
}

export default About;