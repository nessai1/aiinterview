import React from "react";
import './Item.css';
import {Status, StatusBadge} from "@/components/internal/StatusBadge.tsx";
import {AlarmClock} from "lucide-react";
import {ComplexityBadge} from "@/components/internal/ComplexityBadge.tsx";
import {Interview} from "@/lib/interview/interview.ts";


type TProps = {
    interview: Interview
}

const Item: React.FC<TProps> = (props: TProps) => {

    const status = props.interview.complete ? Status.End : Status.Active;
    const statusText = props.interview.complete ? 'Интервью окончено' : 'Интервью в процессе';
    const timer = (props.interview.timing - (props.interview.timing % 60)) / 60;

    return (
        <div className="item">
            <div className="flex-grow">
                <div className="title">
                    {props.interview.title}
                </div>
                <div className="flex items-center mt-1 font-light">
                    <AlarmClock size={15}/>
                    <div className="block ml-1 text-sm">{timer} мин.</div>
                </div>
                <div className="mt-2 flex space-x-1 w-full flex-wrap">
                    {props.interview.sections.map((val, index) => {
                        return <ComplexityBadge complexity={val.grade} key={index}>{val.name}</ComplexityBadge>
                    })}
                </div>
            </div>
            <div>
                <StatusBadge status={status}>{statusText}</StatusBadge>
            </div>
        </div>
    );
}

export default Item;