package service

func (s *UrlService) getUrlRecord(shortUrl string) (string, error) {
	record, err := s.urlRecordDao.Find(shortUrl)
	if err != nil {
		return "", err
	}

	go s.reportStats(shortUrl)

	return record.LongUrl, nil
}

func (s *UrlService) reportStats(shortUrl string) {
	//todo: implement this
}
