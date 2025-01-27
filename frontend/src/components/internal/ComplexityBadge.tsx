import React from "react";
import {Tooltip, TooltipContent, TooltipProvider, TooltipTrigger} from "@/components/ui/tooltip.tsx";
import {Badge} from "@/components/ui/badge.tsx";
import './ComplexityBadge.css';
import {Grade} from "@/lib/interview/interview.ts";

type TProps = {
    complexity: Grade
    children: React.ReactNode
}

const ComplexityBadge: React.FC<TProps> = (props: TProps)=> {
    let hint = '';
    switch (props.complexity) {
        case Grade.Junior:
            hint = 'уровень junior';
            break;
        case Grade.Middle:
            hint = 'уровень middle';
            break;
        case Grade.Senior:
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


export {ComplexityBadge};