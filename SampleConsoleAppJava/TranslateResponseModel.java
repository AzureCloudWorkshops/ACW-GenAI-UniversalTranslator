public class TranslateResponseModel {
    private DetectedLanguage detectedLanguage;
    private Translations[] translations;

    public DetectedLanguage getDetectedLanguage() {
        return detectedLanguage;
    }

    public void setDetectedLanguage(DetectedLanguage detectedLanguage) {
        this.detectedLanguage = detectedLanguage;
    }

    public Translations[] getTranslations() {
        return translations;
    }

    public void setTranslations(Translations[] translations) {
        this.translations = translations;
    }

    public static class DetectedLanguage {
        private String language;
        private double score;

        public String getLanguage() {
            return language;
        }

        public void setLanguage(String language) {
            this.language = language;
        }

        public double getScore() {
            return score;
        }

        public void setScore(double score) {
            this.score = score;
        }
    }

    public static class Translations {
        private String text;
        private String to;

        public String getText() {
            return text;
        }

        public void setText(String text) {
            this.text = text;
        }

        public String getTo() {
            return to;
        }

        public void setTo(String to) {
            this.to = to;
        }
    }
}
