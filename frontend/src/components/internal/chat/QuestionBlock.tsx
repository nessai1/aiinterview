import React from "react";
import './QuestionBlock.css';


type TProps = {
}


const QuestionBlock: React.FC<TProps> = (props: TProps) => {

    return (
        <div className="questionBlock">
            <div className="question">
                <div className="questionHeader">
                    Вопрос 1:
                </div>
                <div className="questionText">
                    <input type="text" className="questionInput" />
                </div>
            </div>
        </div>
    );
}

export default QuestionBlock;