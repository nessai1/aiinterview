import React from "react";
import './QuestionBlock.css';
import MessageEditor from "@/components/internal/chat/editor/MessageEditor.tsx";


type TProps = {
}


const testQuestion = `
    
`;


const QuestionBlock: React.FC<TProps> = (props: TProps) => {

    return (
        <div className="questionBlock">
            <div className="question">
                <div className="questionHeader">
                    Вопрос 1:
                </div>
                <div className="questionText">
                    123
                </div>
                <div className="questionMessage">
                    <MessageEditor />
                </div>
            </div>
        </div>
    );
}

export default QuestionBlock;