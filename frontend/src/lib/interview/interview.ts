type Interview = {
    uuid: string,
    title: string,
    complete: boolean,
    timing: number,
    sections: Section[],
    feedback: string,
    seconds_left: number,
};


type Topic = {
    name: string
    grade: Grade
}

type Section = {
    uuid: string,
    interview_uuid: string,
    name: string,
    complete: boolean,
    color: string,
    questions: Question[]
    position: number,
    grade: Grade,
}

type Question = {
    uuid: string,
    section_uuid: string,
    position: number,
    question: string,
    answer: string,
    feedback: string,
    done: boolean
}

enum Grade {
    Junior = 'junior',
    Middle = 'middle',
    Senior = 'senior'
}

export {
    Grade
};

export type {
    Interview,
    Topic,
    Section,
    Question
}

