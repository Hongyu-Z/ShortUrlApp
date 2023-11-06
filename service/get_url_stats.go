package service

func (s *UrlService) getUrlStats(shortUrl string) (int, int, int, error) {
	return s.UrlStatsDao.GetCount(shortUrl)
}
