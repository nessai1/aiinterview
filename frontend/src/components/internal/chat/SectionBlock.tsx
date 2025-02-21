import React from "react";
import {CountdownTimer} from "@/components/internal/chat/CountdownTimer.tsx";
import {Section} from "@/lib/interview/interview.ts";
import './SectionBlock.css';


type TProps = {
    section: Section
}

const SectionBlock: React.FC<TProps> = (props: TProps) => {

    return (
        <div className="sectionBlock">
            <div className="sectionHeader">
                <div className="circleNumber">1</div>
                <div className="sectionTitle">Алгоритмы</div>
            </div>
        </div>
    );
}

export default SectionBlock;