import { CodeIcon } from "@heroicons/react/solid";
import React from "react";
import { projects } from "../data";

function Projects() {
    return (
        <section id="projects" className="text-gray-400 bg-gray-900 body-font">
            <div className="container px-5 py-10 mx-auto text-center lg:px-40">
                <div className="flex flex-col w-full mb-20">
                    <CodeIcon className="mx-auto inline-block w-10 mb-4" />
                    <h1 className="sm:text-4xl text-3xl font-medium title-font mb-4 text-white">
                        Projects I've Worked On
                    </h1>
                    <p className="lg:w-2/3 mx-auto leading-relaxed text-base">
                        Check some of the projects I have recently been working on. Feel free to use these and any other projects listed on my GitHub.
                    </p>
                </div>
                <div className="flex flex-wrap -m-4">
                    {projects.map((project) => (
                        <a
                            // href={project.link} == Removed to implement buttons for links to GIT and LIVE
                            key={project.image}
                            className="sm:w-1/2 w-100 p-4"
                            target="_blank" rel="noreferrer">
                            <div className="flex relative">
                                <img alt="gallery" className="absolute inset-0 w-full h-full object-cover object-center" src={project.image} />
                                <div className="px-8 py-10 relative z-9 w-full border-4 border-gray-800 bg-gray-900 opacity-0 duration-300 hover:transition-opacity hover:duration-300 hover:opacity-100">
                                    <h2 className="tracking-widest text-sm title-font font-medium text-green-400 mb-1">
                                        {project.subtitle}
                                    </h2>
                                    <h1 className="title-font text-lg font-medium text-white mb-3">
                                        {project.title}
                                    </h1>
                                    <p className="leading-relaxed mb-3">{project.description}</p>
                                    <div class="bg-gray-900 p-4">
                                        <a href={project.git} target="_blank" rel="noreferrer">
                                            <button class="inline-flex items-center px-4 py-2 bg-gray-400 hover:bg-gray-500 text-gray-800 text-sm font-medium rounded-md mr-2">
                                                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="h-5 w-5 mr-2">
                                                    <path stroke-linecap="round" stroke-linejoin="round" d="M6.75 7.5l3 2.25-3 2.25m4.5 0h3m-9 8.25h13.5A2.25 2.25 0 0021 18V6a2.25 2.25 0 00-2.25-2.25H5.25A2.25 2.25 0 003 6v12a2.25 2.25 0 002.25 2.25z" />
                                                </svg>
                                                GitHub
                                            </button>
                                        </a>
                                        <a href={project.link} target="_blank" rel="noreferrer">
                                            <button class="inline-flex items-center px-4 py-2 bg-red-600 hover:bg-red-700 text-white text-sm font-medium rounded-md">
                                                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" class="h-5 w-5 mr-2">
                                                    <path fill-rule="evenodd" d="M4.25 5.5a.75.75 0 00-.75.75v8.5c0 .414.336.75.75.75h8.5a.75.75 0 00.75-.75v-4a.75.75 0 011.5 0v4A2.25 2.25 0 0112.75 17h-8.5A2.25 2.25 0 012 14.75v-8.5A2.25 2.25 0 014.25 4h5a.75.75 0 010 1.5h-5z" clip-rule="evenodd" />
                                                    <path fill-rule="evenodd" d="M6.194 12.753a.75.75 0 001.06.053L16.5 4.44v2.81a.75.75 0 001.5 0v-4.5a.75.75 0 00-.75-.75h-4.5a.75.75 0 000 1.5h2.553l-9.056 8.194a.75.75 0 00-.053 1.06z" clip-rule="evenodd" />
                                                </svg>
                                                Live App
                                            </button>
                                        </a>
                                    </div>
                                </div>
                            </div>
                        </a>
                    ))}
                </div>
            </div>
        </section>
    );
}

export default Projects;