import React from "react";
import {Section} from "@/lib/interview/interview.ts";
import './SectionBlock.css';
import QuestionBlock from "@/components/internal/chat/QuestionBlock.tsx";


type TProps = {
    section: Section,
    interviewComplete: boolean
}

const SectionBlock: React.FC<TProps> = (props: TProps) => {
    return (
        <div className="sectionBlock" style={{backgroundColor: "#"+props.section.color}}>
            <div className="sectionHeader">
                <div className="circleNumber">{props.section.position+1}</div>
                <div className="sectionTitle">{props.section.name}</div>
                <div className={"sectionGrade " + props.section.grade}>{props.section.grade}</div>
            </div>
            <div className="sectionQuestions">
                {
                    props.section.questions.map((question) => (
                        <QuestionBlock key={question.position} question={question} onAnswer={() => {}} showFeedback={props.interviewComplete} />
                    ))
                }
            </div>
        </div>
    );
}

export default SectionBlock;