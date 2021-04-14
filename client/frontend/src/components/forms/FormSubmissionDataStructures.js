class FormFieldValidity {
    constructor(formFieldID, valid){
        this.formFieldID = formFieldID;
        this.valid = valid;
    }
}

class FormSubmissionValidityObject {
    constructor(fieldValidityArray){
        this.fieldValidityArray = fieldValidityArray;
        let overallSubmissionIsValid = true;
        for (let i=0;i<this.fieldValidityArray.length;i++){
            const isValid = this.fieldValidityArray[i].valid;
            if (!isValid){
                overallSubmissionIsValid = false;
            }
        }
        this.validSubmission = overallSubmissionIsValid;
    }
}

export {FormFieldValidity, FormSubmissionValidityObject};