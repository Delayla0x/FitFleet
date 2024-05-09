import axios from 'axios';


interface ClassSchedule {
  id: number;
  name: string;
  date: string;
  time: string;
}


interface BookingDetails {
  classId: number;
  userId: string;
}


const API_BASE_URL: string = process.env.REACT_APP_API_BASE_URL || '';


class ClassService {
  
  static async getClassSchedules(): Promise<ClassSchedule[]> {
    try {
      const response = await axios.get(`${API_BASE_URL}/classes/schedules`);
      return response.data;
    } catch (error) {
      console.error('Error fetching class schedules', error);
      throw error;
    }
  }

  
  static async bookClass(bookingDetails: BookingDetails): Promise<void> {
    try {
      await axios.post(`${API_BASE_URL}/classes/book`, bookingDetails);
      console.log('Class booked successfully');
    } catch (error) {
      console.error('Error booking class', error);
      throw error;
    }
  }

  
  static async cancelReservation(classId: number, userId: string): Promise<void> {
    try {
      await axios.delete(`${API_BASE_URL}/classes/cancel`, {
        data: {
          classId,
          userId,
        },
      });
      console.log('Reservation cancelled successfully');
    } catch (error) {
      console.error('Error cancelling reservation', error);
      throw error;
    }
  }
}


export default ClassService;