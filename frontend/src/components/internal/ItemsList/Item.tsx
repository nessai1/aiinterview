import React from "react";
import './Item.css';
import {Status, StatusBadge} from "@/components/internal/StatusBadge.tsx";
import {AlarmClock} from "lucide-react";
import {Complexity, ComplexityBadge} from "@/components/internal/ComplexityBadge.tsx";


type TProps = {
    title: string

    complete: boolean
}

const Item: React.FC<TProps> = (props: TProps) => {

    const status = props.complete ? Status.End : Status.Active;
    const statusText = props.complete ? 'Интервью окончено' : 'Интервью в процессе';
    return (
        <div className="item">
            <div className="flex-grow">
                <div className="title">
                    Подготовка к собеседованию в Google
                </div>
                <div className="flex items-center mt-1 font-light">
                    <AlarmClock size={15}/>
                    <div className="block ml-1 text-sm">30 минут</div>
                </div>
                <div className="mt-2 flex space-x-1 w-full flex-wrap">
                    <ComplexityBadge complexity={Complexity.Senior}>Алгоритмы</ComplexityBadge>
                    <ComplexityBadge complexity={Complexity.Middle}>PHP</ComplexityBadge>
                    <ComplexityBadge complexity={Complexity.Junior}>Golang</ComplexityBadge>
                    <ComplexityBadge complexity={Complexity.Junior}>Golang</ComplexityBadge>
                    <ComplexityBadge complexity={Complexity.Junior}>Golang</ComplexityBadge>
                    <ComplexityBadge complexity={Complexity.Junior}>Golang</ComplexityBadge>
                    <ComplexityBadge complexity={Complexity.Junior}>Golang</ComplexityBadge>
                    <ComplexityBadge complexity={Complexity.Junior}>Golang</ComplexityBadge>
                    <ComplexityBadge complexity={Complexity.Junior}>Golang</ComplexityBadge>
                    <ComplexityBadge complexity={Complexity.Junior}>Golang</ComplexityBadge>
                    <ComplexityBadge complexity={Complexity.Junior}>Golang</ComplexityBadge>
                    <ComplexityBadge complexity={Complexity.Junior}>Golang</ComplexityBadge>
                    <ComplexityBadge complexity={Complexity.Junior}>Golang</ComplexityBadge>
                    <ComplexityBadge complexity={Complexity.Junior}>Golang</ComplexityBadge>
                    <ComplexityBadge complexity={Complexity.Junior}>Golang</ComplexityBadge>
                    <ComplexityBadge complexity={Complexity.Junior}>Golang</ComplexityBadge>
                    <ComplexityBadge complexity={Complexity.Junior}>Golang</ComplexityBadge>
                    <ComplexityBadge complexity={Complexity.Junior}>Golang</ComplexityBadge>
                    <ComplexityBadge complexity={Complexity.Junior}>Golang</ComplexityBadge>
                </div>
            </div>
            <div>
                <StatusBadge status={status}>{statusText}</StatusBadge>
            </div>
        </div>
    );
}

export default Item;