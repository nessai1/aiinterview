import React from "react";
import './Item.css';
import {Status, StatusBadge} from "@/components/internal/StatusBadge.tsx";

const Item: React.FC = () => {
    return (
        <div className="item">
            <div>
                <div className="title">
                    Подготовка к собеседованию в Google
                </div>
            </div>
            <div>
                <StatusBadge status={Status.End}>Интервью окончено</StatusBadge>
            </div>
        </div>
    );
}

export default Item;