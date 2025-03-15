import React, { FunctionComponent } from 'react';
import './Markdown.css';

interface TProps {
    children: string,
}

const Markdown: FunctionComponent<TProps> = (props) => {
    return (
        <div className="markdown-content" dangerouslySetInnerHTML={{__html: props.children}}></div>
    );
};

export default Markdown;
