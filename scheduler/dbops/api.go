package dbops

import "log"

/*
1. user->api service->delete video
2. api service->scheduler->write video deletion record
3. timer
4. timer->runner->read wvdr->exec->delete video from folder
 */

func AddVideoDeletionRecord(vid string) error {
	stmtIns, err := dbConn.Prepare("insert into video_del_rec (video_id) values (?)")
	if err != nil{
		log.Printf("Insert del video_id err: %v", err)
		return err
	}
	_, err = stmtIns.Exec(vid)
	if err != nil{
		log.Printf("AddVideoDeletionRecord error: %v", err)
		return err
	}
	defer stmtIns.Close()
	log.Printf("AddVideoDeletionRecord success")
	return nil
}