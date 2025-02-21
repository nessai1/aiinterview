import React from "react";
import {Section} from "@/lib/interview/interview.ts";
import './SectionBlock.css';
import QuestionBlock from "@/components/internal/chat/QuestionBlock.tsx";


type TProps = {
    section: Section
}

const SectionBlock: React.FC<TProps> = (props: TProps) => {

    return (
        <div className="sectionBlock">
            <div className="sectionHeader">
                <div className="circleNumber">1</div>
                <div className="sectionTitle">Алгоритмы</div>
                <div className="sectionGrade middle">middle</div>
            </div>
            <div className="sectionQuestions">
                <QuestionBlock />
            </div>
        </div>
    );
}

export default SectionBlock;