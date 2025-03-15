import React, { FunctionComponent} from 'react';
import './IconButton.css'
import {Tooltip, TooltipContent, TooltipProvider, TooltipTrigger} from "@/components/ui/tooltip.tsx";


interface TProps {
    onClick: () => void,
    tooltip: string
    children: React.ReactNode
}


const IconButton: FunctionComponent<TProps> = (props: TProps) => {

    return (
        <TooltipProvider>
            <Tooltip>
                <TooltipTrigger asChild>
                    <button className="icon-button" onClick={props.onClick}>
                        {props.children}
                    </button>
                </TooltipTrigger>
                <TooltipContent>
                    {props.tooltip}
                </TooltipContent>
            </Tooltip>
        </TooltipProvider>
    );
};

export default IconButton;
