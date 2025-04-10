import React, {useEffect, useRef, useState} from "react";
import {CountdownTimer} from "@/components/internal/chat/CountdownTimer.tsx";
import SectionBlock from "@/components/internal/chat/SectionBlock.tsx";
import {Interview, Question, Section} from "@/lib/interview/interview.ts";
import {AxiosError} from "axios";
import {ToastAction} from "@/components/ui/toast.tsx";
import {useToast} from "@/hooks/use-toast.ts";
import {Skeleton} from "@/components/ui/skeleton.tsx";
import Markdown from "@/components/internal/chat/editor/Markdown.tsx";


type TProps = {
    interviewId: string
}

const Chat: React.FC<TProps> = (props: TProps) => {
    const { toast } = useToast();

    const [ interview, setInterview ] = useState<Interview|null>(null)

    const [activeSections, setActiveSections] = useState<Section[]|null>();
    const [interviewComplete, setInterviewComplete] = useState(false);

    const feedbackRef = useRef<HTMLDivElement | null>(null);
    const cdh = (interviewUUID: string) => {
        network.createInterviewFeedback(interviewUUID).then(() => {
            location.reload();
        }).catch((err: AxiosError) => {
            toast({
                title: 'Не могу закончить интервью',
                description: `Ошибка сети: [${err.code}] ${err.message}`,
                action: (
                    <ToastAction altText="Goto schedule to undo" onClick={() => cdh(interviewUUID)}>Повторить</ToastAction>
                ),
                variant: "destructive",
            });
        });
    };

    const retry = () => {
        network.loadInterview(props.interviewId).then((loadedInterview: Interview) => {
            setInterview(loadedInterview);
            const sections: Section[] = [];
            loadedInterview.sections.sort((a, b) => a.position - b.position);
            loadedInterview.sections.forEach((section) => {
                if (section.questions.length > 0)
                {
                    sections.push(section);
                }
            });

            setInterviewComplete(loadedInterview.complete);

            setActiveSections(sections);
        }).catch((err: AxiosError) => {
            toast({
                title: 'Упс! Интервью не загрузилось ;(',
                description: `Ошибка сети: [${err.code}] ${err.message}`,
                action: (
                    <ToastAction altText="Goto schedule to undo" onClick={retry}>Повторить</ToastAction>
                ),
                variant: "destructive",
            });
        });
    };

    useEffect(() => {
        retry();
    }, []);

    useEffect(() => {
        if (interviewComplete) {
            setTimeout(() => {
                feedbackRef.current?.scrollIntoView();
            }, 400)
        }
    }, [interviewComplete]);

    return (
        <>
            {interview === null || activeSections == null ? <Skeleton className="w-full h-[500px] rounded-lg bg-zinc-900" />
                :
                <>
                    {!interview.complete && <CountdownTimer onEnd={() => {
                        const uuid = interview?.uuid;
                        if (uuid)
                        {
                            cdh(uuid);
                        }
                    }} seconds={interview.seconds_left}/>}
                    {activeSections.map((section) => (
                        <SectionBlock key={section.position} section={section} onGetNextSection={(currentPos: number, actualQuestions: Question[]) => {
                            for (let i = 0; i < activeSections.length; i++) {
                                if (activeSections[i].position === currentPos) {
                                    activeSections[i].questions = actualQuestions;
                                }
                            }


                            interview?.sections.forEach((section) => {
                                if (section.position === currentPos+1) {

                                    network.getNextSectionQuestion(section.uuid).then((question: Question) => {
                                        section.questions.push(question);
                                        setActiveSections([...activeSections, section]);
                                    }).catch((err: AxiosError) => {
                                        toast({
                                            title: 'Упс! Вопрос новой секции не загрузился ;(',
                                            description: `Попробуйте перегрузиться. Ошибка сети: [${err.code}] ${err.message}`,
                                            variant: "destructive",
                                        });
                                    })
                                }
                            });

                        }} />
                    ))}
                    {
                        interviewComplete &&
                        <div ref={feedbackRef} className="w-full p-5 bg-zinc-900 rounded-lg mt-4">
                            <p className="text-center text-xl text-zinc-100 font-bold mb-3">Обратная связь по интервью</p>
                            <div>
                                <Markdown>
                                    {interview.feedback}
                                </Markdown>
                            </div>
                        </div>
                    }
                </>
            }

        </>
    );
}

export default Chat;