import React from "react";
import {Tooltip, TooltipContent, TooltipProvider, TooltipTrigger} from "@/components/ui/tooltip.tsx";
import {Badge} from "@/components/ui/badge.tsx";
import './ComplexityBadge.css';


enum Complexity {
    Junior = 'junior',
    Middle = 'middle',
    Senior = 'senior'
}

type TProps = {
    complexity: Complexity
    children: React.ReactNode
}

const ComplexityBadge: React.FC<TProps> = (props: TProps)=> {

    let hint = '';
    switch (props.complexity) {
        case Complexity.Junior:
            hint = 'уровень junior';
            break;
        case Complexity.Middle:
            hint = 'уровень middle';
            break;
        case Complexity.Senior:
            hint = 'уровень senior';
            break;
    }


    return (
        <TooltipProvider>
            <Tooltip>
                <TooltipTrigger asChild>
                   <button className="p-0 border-none m-0" disabled>
                       <Badge className={"rounded-3xl complexity-badge " + props.complexity}>{props.children}</Badge>
                   </button>
                </TooltipTrigger>
                <TooltipContent>
                    {hint}
                </TooltipContent>
            </Tooltip>
        </TooltipProvider>
    );
}


export {ComplexityBadge, Complexity};