import React, {useState} from "react";
import {Question, Section} from "@/lib/interview/interview.ts";
import './SectionBlock.css';
import QuestionBlock from "@/components/internal/chat/QuestionBlock.tsx";
import {Loader2} from "lucide-react";
import {AxiosError} from "axios";
import {useToast} from "@/hooks/use-toast.ts";


type TProps = {
    section: Section,
    onGetNextSection: (currentPos: number, actualQuestions: Question[]) => void,
}

const SectionBlock: React.FC<TProps> = (props: TProps) => {

    const {toast} = useToast();

    const [questions, setQuestions] = useState<Question[]>(props.section.questions);
    const [isLoadNextQuestion, setLoadNextQuestion] = useState<boolean>(false);

    return (
        <div className="sectionBlock" style={{backgroundColor: "#"+props.section.color}}>
            <div className="sectionHeader">
                <div className="circleNumber">{props.section.position+1}</div>
                <div className="sectionTitle">{props.section.name}</div>
                <div className={"sectionGrade " + props.section.grade}>{props.section.grade}</div>
            </div>
            <div className="sectionQuestions">
                {
                    questions.map((question) => (
                        <QuestionBlock
                            key={question.position}
                            question={question}
                            onNextSection={(q: Question) => {
                                for (let i = 0; i < questions.length; i++) {
                                    if (questions[i].uuid === q.uuid) {
                                        questions[i] = q;
                                        setQuestions([...questions]);
                                        break;
                                    }
                                }

                                props.onGetNextSection(props.section.position, questions);
                            }}

                            onNextQuestion={(q: Question) => {
                                for (let i = 0; i < questions.length; i++) {
                                    if (questions[i].uuid === q.uuid) {
                                        questions[i] = q;
                                        setQuestions([...questions]);
                                        break;
                                    }
                                }

                                setLoadNextQuestion(true);

                                network.getNextQuestion(props.section.interview_uuid).then((question: Question) => {
                                    setQuestions([...questions, question]);
                                    setLoadNextQuestion(false);
                                }).catch((err: AxiosError) => {
                                    setLoadNextQuestion(false);
                                    toast({
                                        title: 'Упс! Вопрос не загрузился ;(',
                                        description: `Попробуйте перегрузиться. Ошибка сети: [${err.code}] ${err.message}`,
                                        variant: "destructive",
                                    });
                                })
                            }}
                        />
                    ))
                }
                {
                    isLoadNextQuestion &&
                    <div className="flex justify-center align-center">
                        <div className="mr-2">Загрузка следующего вопроса</div><Loader2 className="animate-spin" />
                    </div>
                }
            </div>
        </div>
    );
}

export default SectionBlock;